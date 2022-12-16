package shorten_test

import (
	"context"
	"testing"

	"github.com/ebriussenex/shrt/internal/model"
	"github.com/ebriussenex/shrt/internal/shorten"
	"github.com/ebriussenex/shrt/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


const url = "github.com/ebriussenex/shrt/"
const nunknownId = "not-exist"
const provId = "hello1"

func TestService_Shorten(t *testing.T) {
	t.Run("gen shortening for a given URL", func(t *testing.T) {
		s := shorten.NewService(shorten_store.NewMongoDbMock())
		in := model.ShortenedReq{
			Url:       url,
			CreatedBy: "",
		}
		shortened, err := s.Shorten(context.Background(), in)
		require.NoError(t, err)

		assert.NotEmpty(t, shortened.Id)
		assert.Equal(t, url, shortened.Url)
		assert.NotZero(t, shortened.CreatedAt)
	})

	t.Run("custom id used if provided", func(t *testing.T) {
		s := shorten.NewService(shorten_store.NewMongoDbMock())
		in := model.ShortenedReq{
			Url:       url,
			CreatedBy: "",
			Id:        provId,
		}
		shortened, err := s.Shorten(context.Background(), in)
		require.NoError(t, err)

		assert.NotEmpty(t, shortened.Id)
		assert.Equal(t, provId, shortened.Id)
	})

	t.Run("if provided id already taken throw error", func(t *testing.T) {
		s := shorten.NewService(shorten_store.NewMongoDbMock())
		in := model.ShortenedReq{
			Url:       url,
			CreatedBy: "",
			Id:        provId,
		}
		_, err := s.Shorten(context.Background(), in)
		require.NoError(t, err)


		_, err = s.Shorten(context.Background(), in)
		require.ErrorIs(t, err, model.ErrIdAlreadyExists)
	})
}

func TestService_Get(t *testing.T) {
	t.Run("get non-existent throws error", func(t *testing.T) {

		s := shorten.NewService(shorten_store.NewMongoDbMock())
		_, err := s.Get(context.Background(), nunknownId)
		require.ErrorIs(t, err, model.ErrNotFound)

	})

	t.Run("get created", func(t *testing.T) {

		s := shorten.NewService(shorten_store.NewMongoDbMock())

		in := model.ShortenedReq{
			Url:       url,
			CreatedBy: "",
			Id:        provId,
		}
		shrt, err := s.Shorten(context.Background(), in)
		require.NoError(t, err)

		shrtg, err := s.Get(context.Background(), provId)
		require.NoError(t, err)

		assert.Equal(t, shrt, shrtg)
	})
}

func TestService_Redirect(t *testing.T) {
	t.Run("redirect increased visit count", func(t *testing.T) {
		s := shorten.NewService(shorten_store.NewMongoDbMock())
		in := model.ShortenedReq{
			Url:       url,
			CreatedBy: "",
			Id:        provId,
		}

		_, err := s.Shorten(context.Background(), in)
		require.NoError(t, err)

		rurl, err := s.Redirect(context.Background(), provId)
		require.NoError(t, err)
		require.NotEmpty(t, rurl)

		shrtg, err := s.Get(context.Background(), provId)
		require.NoError(t, err)
		require.NotZero(t, shrtg.UpdatedAt)

		assert.Equal(t, uint64(1), shrtg.VisitCount)

		_, err = s.Redirect(context.Background(), provId)
		require.NoError(t, err)

		shrtg, err = s.Get(context.Background(), provId)
		require.NoError(t, err)

		assert.Equal(t, uint64(2), shrtg.VisitCount)
	})

	t.Run("redirect non-existent throws error", func(t *testing.T) {
		s := shorten.NewService(shorten_store.NewMongoDbMock())
		url, err := s.Redirect(context.Background(), nunknownId)

		assert.ErrorIs(t, err, model.ErrNotFound)
		assert.Equal(t, url, "")
	})
}
