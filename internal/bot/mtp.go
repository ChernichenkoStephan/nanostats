package bot

import (
	"context"
	"fmt"
	"sync"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message/peer"
	"github.com/gotd/td/telegram/peers"
	mtptg "github.com/gotd/td/tg"

	"go.uber.org/zap"
)

type Key string

const CLIENT_KEY Key = `mtp_client`

// memorySession implements in-memory session storage.
// Goroutine-safe.
type memorySession struct {
	mux  sync.RWMutex
	data []byte
}

// LoadSession loads session from memory.
func (s *memorySession) LoadSession(context.Context) ([]byte, error) {
	if s == nil {
		return nil, session.ErrNotFound
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.data) == 0 {
		return nil, session.ErrNotFound
	}

	cpy := append([]byte(nil), s.data...)

	return cpy, nil
}

// StoreSession stores session to memory.
func (s *memorySession) StoreSession(ctx context.Context, data []byte) error {
	s.mux.Lock()
	s.data = data
	s.mux.Unlock()
	return nil
}

type MTPClient struct {
	// *mtp.Client

	lg      *zap.Logger
	session *memorySession

	token   string
	appID   int
	apiHash string
}

func NewMTP(token string, appID int, apiHash string, log *zap.Logger) *MTPClient {
	ss := &memorySession{}

	return &MTPClient{
		// Client:  mtpClient,
		session: ss,
		lg:      log,
		token:   token,
		appID:   appID,
		apiHash: apiHash,
	}
}

func idsToInputMessageClass(ids []int) []mtptg.InputMessageClass {
	imcs := make([]mtptg.InputMessageClass, 0)
	for _, id := range ids {
		imcs = append(imcs, &mtptg.InputMessageID{ID: id})
	}
	return imcs
}

func (c *MTPClient) withSession(ctx context.Context, f func(ctx context.Context) error) error {

	// No graceful shutdown.
	// ctx := context.TODO()

	mtp := telegram.NewClient(c.appID, c.apiHash, telegram.Options{
		SessionStorage: c.session,
		Logger:         c.lg,
	})

	return mtp.Run(ctx, func(ctx context.Context) error {
		// Checking auth status.
		status, err := mtp.Auth().Status(ctx)
		if err != nil {
			return err
		}
		// Can be already authenticated if we have valid session in
		// session storage.
		if !status.Authorized {
			// Otherwise, perform bot authentication.
			if _, err := mtp.Auth().Bot(ctx, c.token); err != nil {
				return err
			}
		}
		c.lg.Info(`MTP Connected`)
		ctx = context.WithValue(ctx, CLIENT_KEY, mtp)

		err = f(ctx)

		return err
	})
}

func (c MTPClient) getFullMessages(ctx context.Context, username string, ids []int) (mtptg.MessageArray, error) {
	t := ctx.Value(CLIENT_KEY)
	if t == nil {
		return mtptg.MessageArray{}, fmt.Errorf(`client not in context`)
	}
	mtp, ok := t.(*telegram.Client)
	if !ok {
		return mtptg.MessageArray{}, fmt.Errorf("client wrong type")
	}
	if len(ids) == 0 {
		return mtptg.MessageArray{}, fmt.Errorf("no messages ids")
	}
	var messages mtptg.MessageArray

	peerManager := peers.Options{
		Logger: c.lg,
	}.Build(mtp.API())

	p, err := peerManager.ResolveDomain(ctx, username)
	if err != nil {
		return mtptg.MessageArray{}, err
	}

	if inputChannel, ok := peer.ToInputChannel(p.InputPeer()); ok {
		IDs := idsToInputMessageClass(ids)

		req := &mtptg.ChannelsGetMessagesRequest{
			// Channel/supergroup
			Channel: inputChannel,

			// IDs of messages to get
			ID: IDs, // []InputMessageClass
		}

		resp, err := mtp.API().ChannelsGetMessages(ctx, req)
		if err != nil {
			c.lg.Error(fmt.Sprintf("%v", err))
			return mtptg.MessageArray{}, err
		}

		var temp interface{} = resp
		messages = temp.(*mtptg.MessagesChannelMessages).MapMessages().AsMessage()

	} else {
		return mtptg.MessageArray{}, fmt.Errorf("not channel")
	}

	return messages, nil
}
