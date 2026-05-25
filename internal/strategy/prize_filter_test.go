package strategy

import (
	"testing"

	"github.com/Hushengs/Immortality/internal/model"
)

func TestPrizeFilterDecide(t *testing.T) {
	filter := NewPrizeFilter(model.StrategyConfig{
		PrizeWhitelist: []string{"红包", "现金"},
		PrizeBlacklist: []string{"粉丝团"},
		MinJoinSeconds: 5,
		MaxJoinSeconds: 600,
		JoinKeywords:   []string{"参与抽奖", "去参与"},
	})

	t.Run("join when prize and button match", func(t *testing.T) {
		decision := filter.Decide(model.PopupAnalysis{
			PrizeText:        "现金红包",
			CountdownSeconds: 30,
			ButtonText:       "参与抽奖",
		})
		if !decision.ShouldJoin {
			t.Fatalf("expected ShouldJoin to be true, got %+v", decision)
		}
	})

	t.Run("skip blacklisted prize", func(t *testing.T) {
		decision := filter.Decide(model.PopupAnalysis{
			PrizeText:        "粉丝团专属奖励",
			CountdownSeconds: 30,
			ButtonText:       "参与抽奖",
		})
		if !decision.ShouldSkip || !decision.ShouldSwitch {
			t.Fatalf("expected skip + switch, got %+v", decision)
		}
	})
}
