package strategy

import (
	"strings"

	"github.com/Hushengs/Immortality/internal/model"
)

type PrizeFilter struct {
	config model.StrategyConfig
}

func NewPrizeFilter(config model.StrategyConfig) *PrizeFilter {
	return &PrizeFilter{config: config}
}

func (f *PrizeFilter) Decide(analysis model.PopupAnalysis) model.Decision {
	if analysis.CannotJoin {
		return model.Decision{ShouldSkip: true, ShouldSwitch: true, Reason: "button indicates cannot join"}
	}
	if analysis.NeedTask {
		return model.Decision{ShouldSkip: true, ShouldSwitch: true, Reason: "task requirement detected"}
	}
	if analysis.NeedFansClub {
		return model.Decision{ShouldSkip: true, ShouldSwitch: true, Reason: "fans club requirement detected"}
	}
	if analysis.AlreadyJoined {
		return model.Decision{ShouldSkip: true, Reason: "already joined"}
	}
	if analysis.CountdownSeconds > 0 && analysis.CountdownSeconds < f.config.MinJoinSeconds {
		return model.Decision{ShouldSkip: true, ShouldSwitch: true, Reason: "countdown too short"}
	}
	if analysis.CountdownSeconds > f.config.MaxJoinSeconds {
		return model.Decision{ShouldSkip: true, ShouldSwitch: true, Reason: "countdown too long"}
	}
	if blocked := firstKeywordMatch(analysis.PrizeText, f.config.PrizeBlacklist); blocked != "" {
		return model.Decision{ShouldSkip: true, ShouldSwitch: true, Reason: "prize blacklisted: " + blocked}
	}
	if len(f.config.PrizeWhitelist) > 0 && firstKeywordMatch(analysis.PrizeText, f.config.PrizeWhitelist) == "" {
		return model.Decision{ShouldSkip: true, ShouldSwitch: true, Reason: "prize not in whitelist"}
	}
	if joinKeyword := firstKeywordMatch(analysis.ButtonText, f.config.JoinKeywords); joinKeyword != "" {
		return model.Decision{ShouldJoin: true, Reason: "button matched join keyword: " + joinKeyword}
	}
	return model.Decision{ShouldSkip: true, ShouldSwitch: true, Reason: "join button not recognized"}
}

func firstKeywordMatch(text string, keywords []string) string {
	text = strings.ToLower(text)
	for _, keyword := range keywords {
		if keyword != "" && strings.Contains(text, strings.ToLower(keyword)) {
			return keyword
		}
	}
	return ""
}
