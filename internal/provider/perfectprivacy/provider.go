package perfectprivacy

import (
	"math/rand"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/perfectprivacy/updater"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	unzipper common.Unzipper, updaterWarner common.Warner,
) *Provider {
	return &Provider{
		storage:    storage,
		randSource: randSource,
		Fetcher:    updater.New(unzipper, updaterWarner),
	}
}

func (p *Provider) Name() string {
	return providers.Perfectprivacy
}
