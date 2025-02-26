package icons

type IconType int

func (d IconType) String() string {
	return [...]string{"unicode", "emoji", "nerdfont", "ascii"}[d]
}

const (
	Unicode IconType = iota
	Emoji
	NerdFont
	Ascii
)

func GetNextStepIcon(variant IconType) string {
	if variant == Emoji {
		return "⏩"
	}
	if variant == NerdFont {
		return ""
	}
	if variant == Ascii {
		return ">"
	}
	return "↪"
}

func GetNoNextStepIcon(variant IconType) string {
	if variant == Emoji {
		return "🎯"
	}
	if variant == NerdFont {
		return ""
	}
	if variant == Ascii {
		return "#"
	}
	return "◎"
}

func GetBranchIcon(variant IconType) string {
	if variant == Emoji {
		return "🌿"
	}
	if variant == NerdFont {
		return ""
	}
	if variant == Ascii {
		return "Branch"
	}
	return "⎇"
}

func GetLinkIcon(variant IconType) string {
	if variant == Emoji {
		return "🔗  "
	}
	if variant == NerdFont {
		return "  "
	}
	return ""
}

func GetCommitIcon(variant IconType) string {
	if variant == Emoji {
		return "📝  "
	}
	if variant == NerdFont {
		return "  "
	}
	if variant == Ascii {
		return ""

	}
	return "✎  "
}
