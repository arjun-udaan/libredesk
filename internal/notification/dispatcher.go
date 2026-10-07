package notifier

import (
	"fmt"
	"sync/atomic"

	"github.com/abhinavxd/libredesk/internal/notification/channels"
	"github.com/abhinavxd/libredesk/internal/notification/models"
)

type NotificationPreferenceStore interface {
	EnabledChannels(recipientIDs []int, nType models.NotificationType) (map[int][]models.NotificationChannel, error)
}

type Dispatcher struct {
	emailEnabled atomic.Bool
	pipeline     channels.Pipeline
	prefs        NotificationPreferenceStore
}

// DispatcherOpts contains options for creating a new Dispatcher.
type DispatcherOpts struct {
	EmailEnabled bool
	Pipeline     channels.Pipeline
	Prefs        NotificationPreferenceStore
}

// NewDispatcher creates a new notification Dispatcher.
func NewDispatcher(opts DispatcherOpts) *Dispatcher {
	d := &Dispatcher{pipeline: opts.Pipeline, prefs: opts.Prefs}
	d.emailEnabled.Store(opts.EmailEnabled)
	return d
}

func (d *Dispatcher) SetEmailEnabled(enabled bool) {
	d.emailEnabled.Store(enabled)
}

func (d *Dispatcher) Send(n models.Notification) ([]models.DeliveryResult, error) {
	recipientIDs := make([]int, len(n.Recipients))
	for i, recipient := range n.Recipients {
		recipientIDs[i] = recipient.UserID
	}
	enabled, err := d.prefs.EnabledChannels(recipientIDs, n.Type)
	if err != nil {
		return nil, fmt.Errorf("fetching notification preferences: %w", err)
	}

	emailEnabled := d.emailEnabled.Load()
	results := make([]models.DeliveryResult, 0, len(n.Recipients))
	for _, recipient := range n.Recipients {
		prefs := enabled[recipient.UserID]
		if !emailEnabled {
			filtered := make([]models.NotificationChannel, 0, len(prefs))
			for _, channel := range prefs {
				if channel != models.NotificationChannelEmail {
					filtered = append(filtered, channel)
				}
			}
			prefs = filtered
		}
		results = append(results, models.DeliveryResult{
			RecipientID: recipient.UserID,
			Channels: d.pipeline.Send(
				channels.Delivery{Recipient: recipient, Notification: n},
				prefs,
			),
		})
	}
	return results, nil
}
