package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCacher_Cacher(t *testing.T) {
	t.Parallel()

	cacher, _ := NewCacher[string, int](3, 4*time.Hour)

	_ = cacher.Set("k1", 1)
	_ = cacher.Set("k2", 2)
	_ = cacher.Set("k3", 3)
	_ = cacher.Set("k4", 4)
	_ = cacher.Set("k5", 5)

	tests := []struct {
		name    string
		cacher  *Cacher[string, int]
		wantLen int
	}{
		{
			name:    "success",
			cacher:  cacher,
			wantLen: 3,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.wantLen, len(tt.cacher.items))
		})
	}
}

func TestNewCacher(t *testing.T) {
	t.Parallel()

	type args struct {
		maxSize int
		ttl     time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *Cacher[string, int]
		wantErr error
	}{
		{
			name: "zero maxSize",
			args: args{
				maxSize: 0,
				ttl:     4 * time.Hour,
			},
			want:    nil,
			wantErr: ErrInvalidMaxSize,
		},
		{
			name: "negative maxSize",
			args: args{
				maxSize: -1,
				ttl:     4 * time.Hour,
			},
			want:    nil,
			wantErr: ErrInvalidMaxSize,
		},
		{
			name: "zero invalidate time",
			args: args{
				maxSize: 1,
				ttl:     0,
			},
			want:    nil,
			wantErr: ErrInvalidTTL,
		},
		{
			name: "negative invalidate time",
			args: args{
				maxSize: 1,
				ttl:     -4 * time.Hour,
			},
			want:    nil,
			wantErr: ErrInvalidTTL,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewCacher[string, int](tt.args.maxSize, tt.args.ttl)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewCacherWithDebug(t *testing.T) {
	t.Parallel()

	cacher, _ := NewCacher[string, int](3, 4*time.Hour, WithDebug[string, int](zap.NewNop()))

	_ = cacher.Set("k1", 1)
	_ = cacher.Set("k2", 2)
	_ = cacher.Set("k3", 3)

	assert.NotEqual(t, nil, cacher.logger)
}

func TestNewCacherWithUpdateLastUsed(t *testing.T) {
	t.Parallel()

	cacher, _ := NewCacher[string, int](3, 4*time.Hour, WithUpdateLastUsed[string, int]())
	assert.True(t, cacher.updateLastUsed)
}

func TestNewCacherWithThreadSafe(t *testing.T) {
	t.Parallel()

	cacher, _ := NewCacher[string, int](3, 4*time.Hour, WithThreadSafe[string, int]())
	assert.True(t, cacher.threadSafe)
}

func TestCacher_LastUsed(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
	}
	tests := []struct {
		name      string
		cacher    *Cacher[string, int]
		args      args
		want      int
		wantOk    bool
		wantItems map[string]item[int]
	}{
		{
			name: "with update last used option",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     5 * time.Minute,
				now: func() time.Time {
					return time.Date(2022, 10, 25, 13, 53, 0, 0, time.UTC)
				},
				items: map[string]item[int]{
					"k1": {
						value:    1,
						lastUsed: time.Date(2022, 10, 25, 13, 50, 0, 0, time.UTC),
					},
				},
				updateLastUsed: true,
			},
			args: args{
				key: "k1",
			},
			want:   1,
			wantOk: true,
			wantItems: map[string]item[int]{
				"k1": {
					value:    1,
					lastUsed: time.Date(2022, 10, 25, 13, 53, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "without update last used option",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     5 * time.Minute,
				now: func() time.Time {
					return time.Date(2022, 10, 25, 13, 53, 0, 0, time.UTC)
				},
				items: map[string]item[int]{
					"k1": {
						value:    1,
						lastUsed: time.Date(2022, 10, 25, 13, 50, 0, 0, time.UTC),
					},
				},
			},
			args: args{
				key: "k1",
			},
			want:   1,
			wantOk: true,
			wantItems: map[string]item[int]{
				"k1": {
					value:    1,
					lastUsed: time.Date(2022, 10, 25, 13, 50, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotOk := tt.cacher.Get(tt.args.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, gotOk)
			assert.Equal(t, tt.wantItems, tt.cacher.items)
		})
	}
}

func TestCacher_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
	}
	tests := []struct {
		name   string
		cacher *Cacher[string, int]
		args   args
		want   int
		wantOk bool
	}{
		{
			name: "empty key",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     4 * time.Hour,
				now:     time.Now,
				items: map[string]item[int]{
					"k1": {
						value:    1,
						lastUsed: time.Now(),
					},
				},
			},
			args: args{
				key: "",
			},
			want:   0,
			wantOk: false,
		},
		{
			name: "not found",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     4 * time.Hour,
				now:     time.Now,
				items: map[string]item[int]{
					"k1": {
						value:    1,
						lastUsed: time.Now(),
					},
				},
			},
			args: args{
				key: "not_found",
			},
			want:   0,
			wantOk: false,
		},
		{
			name: "deleted by ttl",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     4 * time.Hour,
				now: func() time.Time {
					return time.Date(2022, 10, 25, 13, 50, 0, 0, time.UTC).Add(5 * time.Hour)
				},
				items: map[string]item[int]{
					"k1": {
						value:    1,
						lastUsed: time.Date(2022, 10, 25, 13, 50, 0, 0, time.UTC),
					},
					"k2": {
						value:    2,
						lastUsed: time.Date(2022, 10, 25, 14, 50, 0, 0, time.UTC),
					},
					"k3": {
						value:    3,
						lastUsed: time.Date(2022, 10, 25, 15, 50, 0, 0, time.UTC),
					},
				},
			},
			args: args{
				key: "k1",
			},
			want:   0,
			wantOk: false,
		},
		{
			name: "success",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     4 * time.Hour,
				now:     time.Now,
				items: map[string]item[int]{
					"k1": {
						value:    1,
						lastUsed: time.Now(),
					},
					"k2": {
						value:    2,
						lastUsed: time.Now(),
					},
					"k3": {
						value:    3,
						lastUsed: time.Now(),
					},
				},
			},
			args: args{
				key: "k2",
			},
			want:   2,
			wantOk: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotOk := tt.cacher.Get(tt.args.key)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, gotOk)
		})
	}
}

func TestCacher_Set(t *testing.T) {
	t.Parallel()

	type args struct {
		key   string
		value int
	}
	tests := []struct {
		name      string
		cacher    *Cacher[string, int]
		args      args
		wantErr   error
		wantItems map[string]item[int]
	}{
		{
			name: "empty key",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     4 * time.Hour,
				now:     time.Now,
				items:   make(map[string]item[int], 3),
			},
			args: args{
				key:   "",
				value: 1,
			},
			wantErr:   ErrEmptyKey,
			wantItems: make(map[string]item[int], 3),
		},
		{
			name: "success",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     4 * time.Hour,
				now: func() time.Time {
					return time.Date(2022, 10, 25, 13, 50, 0, 0, time.UTC)
				},
				items: make(map[string]item[int], 3),
			},
			args: args{
				key:   "k1",
				value: 1,
			},
			wantErr: nil,
			wantItems: map[string]item[int]{
				"k1": {
					value:    1,
					lastUsed: time.Date(2022, 10, 25, 13, 50, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cacher.Set(tt.args.key, tt.args.value)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantItems, tt.cacher.items)
		})
	}
}

func TestCacher_clearSpace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cacher    *Cacher[string, int]
		wantItems map[string]item[int]
	}{
		{
			name: "delete last used",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     4 * time.Hour,
				now: func() time.Time {
					return time.Date(2022, 10, 25, 13, 50, 0, 0, time.UTC)
				},
				items: map[string]item[int]{
					"k1": {
						value:    1,
						lastUsed: time.Date(2022, 10, 25, 12, 51, 0, 0, time.UTC),
					},
					"k2": {
						value:    2,
						lastUsed: time.Date(2022, 10, 25, 12, 52, 0, 0, time.UTC),
					},
					"k3": {
						value:    3,
						lastUsed: time.Date(2022, 10, 25, 12, 53, 0, 0, time.UTC),
					},
				},
			},
			wantItems: map[string]item[int]{
				"k2": {
					value:    2,
					lastUsed: time.Date(2022, 10, 25, 12, 52, 0, 0, time.UTC),
				},
				"k3": {
					value:    3,
					lastUsed: time.Date(2022, 10, 25, 12, 53, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "delete invalidate records",
			cacher: &Cacher[string, int]{
				maxSize: 3,
				ttl:     4 * time.Hour,
				now: func() time.Time {
					return time.Date(2022, 10, 25, 16, 50, 0, 0, time.UTC)
				},
				items: map[string]item[int]{
					"k1": {
						value:    1,
						lastUsed: time.Date(2022, 10, 25, 12, 40, 0, 0, time.UTC),
					},
					"k2": {
						value:    2,
						lastUsed: time.Date(2022, 10, 25, 12, 45, 0, 0, time.UTC),
					},
					"k3": {
						value:    3,
						lastUsed: time.Date(2022, 10, 25, 15, 50, 0, 0, time.UTC),
					},
				},
			},
			wantItems: map[string]item[int]{
				"k3": {
					value:    3,
					lastUsed: time.Date(2022, 10, 25, 15, 50, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.cacher.clearSpace()
			assert.Equal(t, tt.wantItems, tt.cacher.items)
		})
	}
}

func Test_isEmpty(t *testing.T) {
	t.Parallel()

	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{
				v: "",
			},
			want: true,
		},
		{
			name: "not empty",
			args: args{
				v: "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isEmpty(tt.args.v)
			assert.Equal(t, tt.want, got)
		})
	}
}
