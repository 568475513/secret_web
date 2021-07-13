package enums

const (
	// 付费类型：1-免费、2-单笔、3-付费产品包（包括大专栏）、4-团购、5-单笔的购买赠送、6-产品包的购买赠送、7-问答提问、8-问答偷听、9-购买会员、10-会员的购买赠送、11、活动付费报名 12、打赏  13、拼团单个资源 14、拼团产品包
	PaymentTypeFree           = 1
	PaymentTypeSingle         = 2
	PaymentTypeProduct        = 3
	PaymentTypeGroup          = 4
	PaymentTypeGiftSingle     = 5
	PaymentTypeGiftProduct    = 6
	PaymentTypeQueAsk         = 7
	PaymentTypeQueListen      = 8
	PaymentTypeMember         = 9
	PaymentTypeGiftMember     = 10
	PaymentTypeActivity       = 11
	PaymentTypeReward         = 12
	PaymentTypeTopic          = 3
	PaymentTypeTeamBuySingle  = 13
	PaymentTypeTeamBuyProduct = 14
	PaymentTypeSvip           = 15

	// 资源类型
	ResourceTypeImageText = 1
	ResourceTypeAudio     = 2
	ResourceTypeVideo     = 3
	ResourceTypeLive      = 4

	AGENT_TYPE_MP          = 0 //微信小程序
	AGENT_TYPE_H5          = 1 //微信浏览器
	AGENT_TYPE_WW          = 3 //企业微信浏览器
	AGENT_TYPE_APP         = 4 //小鹅通App
	AGENT_TYPE_XIAOE_AGENT = 5 //小鹅通内嵌浏览器
	AGENT_TYPE_MOBILE      = 6 //手机浏览器
	AGENT_TYPE_PC          = 7 //PC端浏览器
	AGENT_TYPE_TRAINING    = 8 //企业内训

	// 店铺版本常量
	VERSION_TYPE_PROBATION        = 0   //试用版
	VERSION_TYPE_STANDARD         = 4   //标准版版,
	VERSION_TYPE_ONLINE_EDUCATION = 7   //在线教育
	VERSION_TYPE_ADVANCED         = 8   //旗舰版
	VERSION_TYPE_TRAINING_TRY     = 170 // 企学院试用版
	VERSION_TYPE_TRAINING_STD     = 171 // 企学院正式版
	VERSION_TYPE_QLIVE            = 301 // 企业直播
	VERSION_TYPE_ELIVE            = 310 // 鹅直播个人店铺

	// 直播删除状态
	AliveStateDelete = 2

	// 直播类型
	AliveTypeVideo   = 1
	AliveTypePush    = 2
	AliveOldTypePush = 4

	//回放过期类型
	LookBackExpireTypeNever = 1 //永久
	LookBackExpireTypeFixed = 2 //固定日期

	//直播房间页
	AliveRoomPage = 0 // 直播房间页
)
