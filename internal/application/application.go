package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/x-foby/kakafka/internal/kafka"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Application struct {
	ctx context.Context
	cfg Config

	conns map[string]*kafka.Conn
}

func New() *Application {
	return &Application{
		conns: make(map[string]*kafka.Conn),
	}
}

func (a *Application) Startup(ctx context.Context) {
	a.ctx = ctx

	cfg, err := loadConfig(a.ctx)
	if err != nil {
		runtime.LogFatalf(a.ctx, "failed to load config: %s", err.Error())
	}

	a.cfg = *cfg
}

func (a *Application) GetConfigs() Config {
	return a.cfg
}

func (a *Application) CreateProfile(p Profile) error {
	err := a.createProfile(p)
	if err != nil {
		runtime.MessageDialog(
			a.ctx,
			runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: err.Error(),
			},
		)
	}

	return err
}

func (a *Application) createProfile(p Profile) error {
	for _, profile := range a.cfg.Profiles {
		if profile.Name == p.Name {
			return fmt.Errorf("profile %q already exists", p.Name)
		}
	}

	if len(p.Brokers) == 0 {
		return errors.New("empty broker list")
	}

	a.cfg.Profiles = append(a.cfg.Profiles, p)

	return saveConfig(a.ctx, a.cfg)
}

func (a *Application) DeleteProfile(profileName string) error {
	err := a.deleteProfile(profileName)
	if err != nil {
		runtime.MessageDialog(
			a.ctx,
			runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: err.Error(),
			},
		)
	}

	return err
}

func (a *Application) deleteProfile(profileName string) error {
	var profiles []Profile

	for _, profile := range a.cfg.Profiles {
		if profile.Name != profileName {
			profiles = append(profiles, profile)
		}
	}

	a.cfg.Profiles = profiles

	return saveConfig(a.ctx, a.cfg)
}

func (a *Application) Connect(profileName string) error {
	err := a.connect(profileName)
	if err != nil {
		runtime.MessageDialog(
			a.ctx,
			runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: err.Error(),
			},
		)
	}

	return err
}

func (a *Application) connect(profileName string) error {
	runtime.LogDebugf(a.ctx, "lookup for opened connects: %s", profileName)

	if _, ok := a.conns[profileName]; ok {
		return nil
	}

	runtime.LogDebugf(a.ctx, "lookup for profiles: %s", profileName)

	profile, ok := a.getProfileByName(profileName)
	if !ok {
		return fmt.Errorf("profile %q not found", profileName)
	}

	runtime.LogDebugf(a.ctx, "connecting for profiles: %s", profileName)

	conn, err := kafka.Connect(a.ctx, profile.Brokers...)
	if err != nil {
		return err
	}

	runtime.LogDebug(a.ctx, "successifully connected")

	a.conns[profileName] = conn

	return nil
}

func (a *Application) GetTopics(profile string, refresh bool) ([]kafka.Topic, error) {
	topics, err := a.getTopics(profile, refresh)
	if err != nil {
		runtime.MessageDialog(
			a.ctx,
			runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: err.Error(),
			},
		)
	}

	return topics, err
}

func (a *Application) getTopics(profile string, refresh bool) ([]kafka.Topic, error) {
	runtime.LogDebugf(a.ctx, "lookup for opened connects: %s", profile)

	conn, ok := a.conns[profile]
	if !ok {
		return nil, fmt.Errorf("profile %q not connected", profile)
	}

	return conn.GetTopics(a.ctx, refresh)
}

func (a *Application) CreateTopic(profile string, topic kafka.TopicConfig) (*kafka.Topic, error) {
	topics, err := a.createTopic(profile, topic)
	if err != nil {
		runtime.MessageDialog(
			a.ctx,
			runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: err.Error(),
			},
		)
	}

	return topics, err
}

func (a *Application) createTopic(profile string, topic kafka.TopicConfig) (*kafka.Topic, error) {
	runtime.LogDebugf(a.ctx, "lookup for opened connects: %s", profile)

	conn, ok := a.conns[profile]
	if !ok {
		return nil, fmt.Errorf("profile %q not connected", profile)
	}

	runtime.LogDebugf(a.ctx, "creating topic %s with %d partitions and %d replication factor", topic.Topic, topic.NumPartitions, topic.ReplicationFactor)

	return conn.CreateTopic(a.ctx, topic)
}

func (a *Application) DeleteTopic(profile, topic string) error {
	if err := a.deleteTopic(profile, topic); err != nil {
		runtime.MessageDialog(
			a.ctx,
			runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: err.Error(),
			},
		)
	}

	return nil
}

func (a *Application) deleteTopic(profile, topicName string) error {
	runtime.LogDebugf(a.ctx, "lookup for opened connects: %s", profile)

	conn, ok := a.conns[profile]
	if !ok {
		return fmt.Errorf("profile %q not connected", profile)
	}

	runtime.LogDebugf(a.ctx, "deleting topic %s", topicName)

	return conn.DeleteTopic(a.ctx, topicName)
}

func (a *Application) ConsumerOffsets(profile, topic string) ([]kafka.ConsumerOffset, error) {
	offsets, err := a.consumerOffsets(profile, topic)
	if err != nil {
		runtime.MessageDialog(
			a.ctx,
			runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: err.Error(),
			},
		)
	}

	return offsets, err
}

func (a *Application) consumerOffsets(profile, topic string) ([]kafka.ConsumerOffset, error) {
	runtime.LogDebugf(a.ctx, "lookup for opened connects: %s", profile)

	conn, ok := a.conns[profile]
	if !ok {
		return nil, fmt.Errorf("profile %q not connected", profile)
	}

	runtime.LogDebugf(a.ctx, "fetch consumer offsets for topic %s", topic)

	return conn.ConsumerOffsets(a.ctx, topic)
}

func (a *Application) GetMessages(profile, topic, query string) ([]kafka.Message, error) {
	messages, err := a.getMessages(profile, topic, query)
	if err != nil {
		runtime.MessageDialog(
			a.ctx,
			runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: err.Error(),
			},
		)
	}

	return messages, err
}

func (a *Application) getMessages(profile, topic, query string) ([]kafka.Message, error) {
	runtime.LogDebugf(a.ctx, "lookup for opened connects: %s", profile)

	conn, ok := a.conns[profile]
	if !ok {
		return nil, fmt.Errorf("profile %q not connected", profile)
	}

	return conn.GetMessages(a.ctx, topic, query)
}

func (a *Application) getProfileByName(profileName string) (Profile, bool) {
	for _, p := range a.cfg.Profiles {
		if p.Name == profileName {
			return p, true
		}
	}

	return Profile{}, false
}
