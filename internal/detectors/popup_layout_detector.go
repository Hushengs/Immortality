package detectors

import "github.com/Hushengs/Immortality/internal/model"

type PopupLayoutDetector struct {
	config model.RuntimeConfig
}

func NewPopupLayoutDetector(config model.RuntimeConfig) *PopupLayoutDetector {
	return &PopupLayoutDetector{config: config}
}

func (d *PopupLayoutDetector) Detect() model.PopupLayout {
	return model.PopupLayout{
		Name:            "default",
		PrizeRegion:     d.config.ScreenRegions.PopupPrize,
		CountdownRegion: d.config.ScreenRegions.PopupCountdown,
		ButtonRegion:    d.config.ScreenRegions.PopupButton,
	}
}
