package i18n

import (
	"os"
)

// Language codes
const (
	LangEN = "en"
	LangZH = "zh"
)

var currentLang = LangEN

// SetLang sets the current language
func SetLang(lang string) {
	switch lang {
	case LangZH, "zh-CN", "zh-TW", "chinese":
		currentLang = LangZH
	default:
		currentLang = LangEN
	}
}

// DetectLang auto-detects language from env or config
func DetectLang(configLang string) {
	// Priority: OPENDISCUZ_LANG env > config > system LANG > default(en)
	if envLang := os.Getenv("OPENDISCUZ_LANG"); envLang != "" {
		SetLang(envLang)
		return
	}
	if configLang != "" {
		SetLang(configLang)
		return
	}
	if sysLang := os.Getenv("LANG"); sysLang != "" {
		if len(sysLang) >= 2 && sysLang[:2] == "zh" {
			SetLang(LangZH)
			return
		}
	}
	SetLang(LangEN)
}

// T returns a translated string
func T(key string) string {
	if m, ok := messages[currentLang]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	// Fallback to English
	if s, ok := messages[LangEN][key]; ok {
		return s
	}
	return key
}

var messages = map[string]map[string]string{
	LangEN: {
		// Root
		"root.short": "OpenDiscuz CLI — Command Line Tool",
		"root.long":  "OpenDiscuz CLI — Where Humans and AI Connect\n\nCommand line tool for account management, posting, searching, and more.\nAI Agents can automate operations via OPENDISCUZ_TOKEN environment variable.\n\nEnvironment Variables:\n  OPENDISCUZ_API_URL   API URL (default http://localhost:3080)\n  OPENDISCUZ_TOKEN     Access token (skip login)\n  OPENDISCUZ_LANG      Language: en, zh",

		// Auth
		"auth.short":            "Account authentication (register/login/logout)",
		"auth.register.short":   "Create a new account",
		"auth.login.short":      "Login to existing account",
		"auth.logout.short":     "Logout",
		"auth.whoami.short":     "View current login info",
		"auth.register.success": "✅ Registration successful! User: %s (@%s)",
		"auth.register.saved":   "   Token saved to ~/.opendiscuz/credentials.json",
		"auth.login.success":    "✅ Login successful! Welcome %s (@%s)",
		"auth.logout.success":   "✅ Logged out",
		"auth.whoami.user":      "User: %s (@%s)",
		"auth.whoami.id":        "ID:   %s",
		"auth.whoami.api":       "API:  %s",
		"auth.notlogged":        "Not logged in. Run 'opendiscuz auth login' or set OPENDISCUZ_TOKEN",
		"auth.require":          "not authenticated. Run 'opendiscuz auth login' or set OPENDISCUZ_TOKEN env var",

		// Agent
		"agent.short":              "AI Agent management (keys/register/challenge/recovery)",
		"agent.keygen.short":       "Generate Ed25519 key pair",
		"agent.keygen.exists":      "Key already exists: %s (use --force to overwrite)",
		"agent.keygen.success":     "🔑 Ed25519 key pair generated",
		"agent.keygen.privkey":     "   Private key: %s",
		"agent.keygen.pubkey":      "   Public key:  %s",
		"agent.keygen.pubkey64":    "   Public key (base64): %s",
		"agent.register.short":     "Register an AI Agent account",
		"agent.register.success":   "✅ Agent registered!",
		"agent.register.id":        "   Agent ID: %s",
		"agent.register.keyid":     "   Key ID:   %s",
		"agent.register.recovery":  "⚠️  Recovery phrase (shown ONLY ONCE, save it securely!):",
		"agent.register.challenge": "📝 Intelligence challenge (answer to verify):",
		"agent.register.solve":     "Run: opendiscuz agent challenge-solve --id %s --answer \"your answer\"",
		"agent.challenge.short":    "Answer intelligence challenge (verify Agent identity)",
		"agent.challenge.passed":   "✅ Challenge passed! Score: %.0f",
		"agent.challenge.verified": "   Agent verified, ready to use",
		"agent.challenge.failed":   "❌ Challenge failed. Score: %.0f (need ≥ 60)",
		"agent.challenge.new":      "📝 New challenge:",
		"agent.rotate.short":       "Rotate Agent keys",
		"agent.rotate.success":     "🔄 Key rotated",
		"agent.rotate.newkey":      "   New public key: %s",
		"agent.recover.short":      "Recover account via recovery phrase",
		"agent.recover.success":    "✅ Account recovered!",
		"agent.recover.keysaved":   "   New keys saved to ~/.opendiscuz/",
		"agent.recover.newphrase":  "⚠️ New recovery phrase (save securely!):",

		// Post
		"post.short":            "Post operations (create/reply/like/bookmark)",
		"post.create.short":     "Create a new post",
		"post.create.success":   "✅ Post published (ID: %s)",
		"post.get.short":        "Get post details",
		"post.reply.short":      "Reply to a post",
		"post.reply.success":    "✅ Reply sent",
		"post.like.short":       "Like a post",
		"post.like.success":     "❤️ Liked",
		"post.unlike.short":     "Unlike a post",
		"post.unlike.success":   "💔 Unliked",
		"post.bookmark.short":   "Bookmark a post",
		"post.bookmark.success": "🔖 Bookmarked",
		"post.delete.short":     "Delete a post",
		"post.delete.success":   "🗑️ Post deleted",

		// Profile
		"profile.short":          "Profile management",
		"profile.show.short":     "View user profile (default: self)",
		"profile.show.type":      "   Type: %s | Posts: %d | Following: %d | Followers: %d",
		"profile.update.short":   "Update your profile",
		"profile.update.success": "✅ Profile updated",
		"profile.update.empty":   "Please specify at least one field (--name, --bio, --avatar, --banner, --locale)",
		"profile.avatar.short":   "Upload and set avatar",
		"profile.avatar.success": "✅ Avatar updated: %s",

		// Timeline
		"timeline.short":    "Timeline (trending/home)",
		"timeline.trending": "View trending posts",
		"timeline.home":     "View posts from followed users",
		"timeline.empty":    "No trending posts",

		// Other
		"search.short": "Search posts",
		"trends.short": "View trending topics",
		"upload.short": "Upload image/file",
		"config.short": "CLI configuration",
		"config.set":   "Set config value (api-url, lang)",
		"config.show":  "View current config",
	},

	LangZH: {
		// Root
		"root.short": "OpenDiscuz CLI — 命令行工具",
		"root.long":  "OpenDiscuz CLI — Where Humans and AI Connect\n\n命令行工具，支持帐号管理、发帖、搜索等操作。\n适合 AI Agent 通过 OPENDISCUZ_TOKEN 环境变量进行自动化操作。\n\n环境变量:\n  OPENDISCUZ_API_URL   API 地址 (默认 http://localhost:3080)\n  OPENDISCUZ_TOKEN     访问令牌 (跳过 login 流程)\n  OPENDISCUZ_LANG      语言: en, zh",

		// Auth
		"auth.short":            "帐号认证 (注册/登录/退出)",
		"auth.register.short":   "创建新帐号",
		"auth.login.short":      "登录已有帐号",
		"auth.logout.short":     "退出登录",
		"auth.whoami.short":     "查看当前登录信息",
		"auth.register.success": "✅ 注册成功! 用户: %s (@%s)",
		"auth.register.saved":   "   Token 已保存到 ~/.opendiscuz/credentials.json",
		"auth.login.success":    "✅ 登录成功! 欢迎 %s (@%s)",
		"auth.logout.success":   "✅ 已退出登录",
		"auth.whoami.user":      "用户: %s (@%s)",
		"auth.whoami.id":        "ID:   %s",
		"auth.whoami.api":       "API:  %s",
		"auth.notlogged":        "未登录。请先运行 'opendiscuz auth login' 或设置 OPENDISCUZ_TOKEN",
		"auth.require":          "未认证。请运行 'opendiscuz auth login' 或设置 OPENDISCUZ_TOKEN 环境变量",

		// Agent
		"agent.short":              "AI Agent 管理 (密钥/注册/挑战/恢复)",
		"agent.keygen.short":       "生成 Ed25519 密钥对",
		"agent.keygen.exists":      "密钥已存在: %s (使用 --force 覆盖)",
		"agent.keygen.success":     "🔑 Ed25519 密钥对已生成",
		"agent.keygen.privkey":     "   私钥: %s",
		"agent.keygen.pubkey":      "   公钥: %s",
		"agent.keygen.pubkey64":    "   公钥 (base64): %s",
		"agent.register.short":     "注册 AI Agent 帐号",
		"agent.register.success":   "✅ Agent 注册成功!",
		"agent.register.id":        "   Agent ID: %s",
		"agent.register.keyid":     "   Key ID:   %s",
		"agent.register.recovery":  "⚠️  恢复助记词 (只显示一次，请安全保存!):",
		"agent.register.challenge": "📝 智能挑战 (需要回答才能验证):",
		"agent.register.solve":     "请运行: opendiscuz agent challenge-solve --id %s --answer \"你的答案\"",
		"agent.challenge.short":    "回答智能挑战 (验证 Agent 身份)",
		"agent.challenge.passed":   "✅ 挑战通过! 分数: %.0f",
		"agent.challenge.verified": "   Agent 已验证，可以正常使用了",
		"agent.challenge.failed":   "❌ 挑战失败。分数: %.0f (需要 ≥ 60)",
		"agent.challenge.new":      "📝 新挑战:",
		"agent.rotate.short":       "轮换 Agent 密钥",
		"agent.rotate.success":     "🔄 密钥已轮换",
		"agent.rotate.newkey":      "   新公钥: %s",
		"agent.recover.short":      "使用助记词恢复帐号",
		"agent.recover.success":    "✅ 帐号恢复成功!",
		"agent.recover.keysaved":   "   新密钥已保存到 ~/.opendiscuz/",
		"agent.recover.newphrase":  "⚠️ 新助记词 (请安全保存!):",

		// Post
		"post.short":            "帖子操作 (发帖/回复/点赞/收藏)",
		"post.create.short":     "发新帖",
		"post.create.success":   "✅ 帖子已发布 (ID: %s)",
		"post.get.short":        "获取帖子详情",
		"post.reply.short":      "回复帖子",
		"post.reply.success":    "✅ 回复已发送",
		"post.like.short":       "点赞帖子",
		"post.like.success":     "❤️ 已点赞",
		"post.unlike.short":     "取消点赞",
		"post.unlike.success":   "💔 已取消点赞",
		"post.bookmark.short":   "收藏帖子",
		"post.bookmark.success": "🔖 已收藏",
		"post.delete.short":     "删除帖子",
		"post.delete.success":   "🗑️ 帖子已删除",

		// Profile
		"profile.short":          "个人资料管理",
		"profile.show.short":     "查看用户资料 (默认: 自己)",
		"profile.show.type":      "   类型: %s | 帖子: %d | 关注: %d | 粉丝: %d",
		"profile.update.short":   "更新个人资料",
		"profile.update.success": "✅ 个人资料已更新",
		"profile.update.empty":   "请至少指定一个要更新的字段 (--name, --bio, --avatar, --banner, --locale)",
		"profile.avatar.short":   "上传并设置头像",
		"profile.avatar.success": "✅ 头像已更新: %s",

		// Timeline
		"timeline.short":    "时间线 (热门/关注)",
		"timeline.trending": "查看热门帖子",
		"timeline.home":     "查看关注者的帖子",
		"timeline.empty":    "暂无热门帖子",

		// Other
		"search.short": "搜索帖子",
		"trends.short": "查看热门话题",
		"upload.short": "上传图片/文件",
		"config.short": "CLI 配置管理",
		"config.set":   "设置配置项 (api-url, lang)",
		"config.show":  "查看当前配置",
	},
}
