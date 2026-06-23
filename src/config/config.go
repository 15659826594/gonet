package config

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/spf13/viper"
)

type Config struct {
	// 后台入口
	AdminName string `mapstructure:"admin_name"`
	// 应用设置相关配置
	AppNamespace        string   `mapstructure:"app_namespace"`
	AppDebug            bool     `mapstructure:"app_debug"`
	AppTrace            bool     `mapstructure:"app_trace"`
	AppStatus           string   `mapstructure:"app_status"`
	AppMultiModule      bool     `mapstructure:"app_multi_module"`
	AutoBindModule      bool     `mapstructure:"auto_bind_module"`
	RootNamespace       []string `mapstructure:"root_namespace"`
	ExtraFileList       []string `mapstructure:"extra_file_list"`
	DefaultReturnType   string   `mapstructure:"default_return_type"`
	DefaultAjaxReturn   string   `mapstructure:"default_ajax_return"`
	DefaultJsonpHandler string   `mapstructure:"default_jsonp_handler"`
	VarJsonpHandler     string   `mapstructure:"var_jsonp_handler"`
	DefaultTimezone     string   `mapstructure:"default_timezone"`
	LangSwitchOn        bool     `mapstructure:"lang_switch_on"`
	DefaultFilter       string   `mapstructure:"default_filter"`
	DefaultLang         string   `mapstructure:"default_lang"`
	AllowLangList       []string `mapstructure:"allow_lang_list"`
	ClassSuffix         bool     `mapstructure:"class_suffix"`
	ControllerSuffix    bool     `mapstructure:"controller_suffix"`
	HttpAgentIp         string   `mapstructure:"http_agent_ip"`

	// 模块设置相关配置
	DefaultModule        string   `mapstructure:"default_module"`
	DenyModuleList       []string `mapstructure:"deny_module_list"`
	DefaultController    string   `mapstructure:"default_controller"`
	DefaultAction        string   `mapstructure:"default_action"`
	DefaultValidate      string   `mapstructure:"default_validate"`
	EmptyController      string   `mapstructure:"empty_controller"`
	ActionSuffix         string   `mapstructure:"action_suffix"`
	ControllerAutoSearch bool     `mapstructure:"controller_auto_search"`

	// URL设置相关配置
	VarPathinfo        string   `mapstructure:"var_pathinfo"`
	PathinfoFetch      []string `mapstructure:"pathinfo_fetch"`
	PathinfoDepr       string   `mapstructure:"pathinfo_depr"`
	UrlHtmlSuffix      string   `mapstructure:"url_html_suffix"`
	UrlCommonParam     bool     `mapstructure:"url_common_param"`
	UrlParamType       int      `mapstructure:"url_param_type"`
	UrlRouteOn         bool     `mapstructure:"url_route_on"`
	RouteCompleteMatch bool     `mapstructure:"route_complete_match"`
	RouteConfigFile    []string `mapstructure:"route_config_file"`
	UrlRouteMust       bool     `mapstructure:"url_route_must"`
	UrlDomainDeploy    bool     `mapstructure:"url_domain_deploy"`
	UrlDomainRoot      string   `mapstructure:"url_domain_root"`
	UrlConvert         bool     `mapstructure:"url_convert"`
	UrlControllerLayer string   `mapstructure:"url_controller_layer"`
	VarMethod          string   `mapstructure:"var_method"`
	VarAjax            string   `mapstructure:"var_ajax"`
	VarPjax            string   `mapstructure:"var_pjax"`
	RequestCache       bool     `mapstructure:"request_cache"`
	RequestCacheExpire int      `mapstructure:"request_cache_expire"`

	// 在 Config 结构体中追加以下字段
	Template              *Template              `mapstructure:"template"`
	ViewReplaceStr        map[string]string      `mapstructure:"view_replace_str"`
	DispatchSuccessTmpl   string                 `mapstructure:"dispatch_success_tmpl"`
	DispatchErrorTmpl     string                 `mapstructure:"dispatch_error_tmpl"`
	ExceptionTmpl         string                 `mapstructure:"exception_tmpl"`
	HttpExceptionTemplate map[string]interface{} `mapstructure:"http_exception_template"`
	ErrorMessage          string                 `mapstructure:"error_message"`
	ShowErrorMsg          bool                   `mapstructure:"show_error_msg"`
	Log                   *Log                   `mapstructure:"log"`
	Trace                 *Trace                 `mapstructure:"trace"`
	Cache                 *Cache                 `mapstructure:"cache"`
	Session               *Session               `mapstructure:"session"`
	Cookie                *Cookie                `mapstructure:"cookie"`
	Paginate              *Paginate              `mapstructure:"paginate"`
	Captcha               *Captcha               `mapstructure:"captcha"`
	Token                 *Token                 `mapstructure:"token"`
	Fastadmin             *Fastadmin             `mapstructure:"fastadmin"`
}

func (c *Config) Mode() string {
	if viper.IsSet("APP_ENV") && slices.Contains([]string{gin.ReleaseMode, gin.DebugMode, gin.TestMode}, viper.GetString("APP_ENV")) {
		return viper.GetString("APP_ENV")
	}
	if c.AppDebug {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}

func (c *Config) ExceptionHandle() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

type Template struct {
	Type        string   `mapstructure:"type"`
	ViewPath    string   `mapstructure:"view_path"`
	ViewSuffix  []string `mapstructure:"view_suffix"`
	ViewDepr    string   `mapstructure:"view_depr"`
	TplBegin    string   `mapstructure:"tpl_begin"`
	TplEnd      string   `mapstructure:"tpl_end"`
	TaglibBegin string   `mapstructure:"taglib_begin"`
	TaglibEnd   string   `mapstructure:"taglib_end"`
	TplCache    bool     `mapstructure:"tpl_cache"`
}

func (t *Template) Delims() render.Delims {
	delims := render.Delims{
		Left:  t.TplBegin,
		Right: t.TplEnd,
	}
	if delims.Left == "" {
		delims.Left = "{{"
	}
	if delims.Right == "" {
		delims.Right = "}}"
	}
	return delims
}

type Log struct {
	Prefix      string   `mapstructure:"prefix"`
	Type        string   `mapstructure:"type"`
	Path        string   `mapstructure:"path"`
	Level       string   `mapstructure:"level"`
	FileSize    int      `mapstructure:"file_size"`
	ExcludeUri  []string `mapstructure:"exclude_uri"`
	DateFormat  string   `mapstructure:"date_format"`
	JsonFormat  bool     `mapstructure:"json_format"`
	ApmEnabled  bool     `mapstructure:"apm_enabled"`
	RotateByDay bool     `mapstructure:"rotate_by_day"`
}

type Trace struct {
	Type string `mapstructure:"type"`
}

type Cache struct {
	Type   string `mapstructure:"type"`
	Path   string `mapstructure:"path"`
	Prefix string `mapstructure:"prefix"`
	Expire int    `mapstructure:"expire"`
}

type Session struct {
	Id           string `mapstructure:"id"`
	VarSessionId string `mapstructure:"var_session_id"`
	Prefix       string `mapstructure:"prefix"`
	Type         string `mapstructure:"type"`
	AutoStart    bool   `mapstructure:"auto_start"`
}

type Cookie struct {
	Prefix    string `mapstructure:"prefix"`
	Expire    int    `mapstructure:"expire"`
	Path      string `mapstructure:"path"`
	Domain    string `mapstructure:"domain"`
	Secure    bool   `mapstructure:"secure"`
	Httponly  string `mapstructure:"httponly"`
	Setcookie bool   `mapstructure:"setcookie"`
}

type Paginate struct {
	Type     string `mapstructure:"type"`
	VarPage  string `mapstructure:"var_page"`
	ListRows int    `mapstructure:"list_rows"`
}

type Captcha struct {
	CodeSet  string `mapstructure:"code_set"`
	FontSize int    `mapstructure:"font_size"`
	UseCurve bool   `mapstructure:"use_curve"`
	UseZh    bool   `mapstructure:"use_zh"`
	ImageH   int    `mapstructure:"image_h"`
	ImageW   int    `mapstructure:"image_w"`
	Length   int    `mapstructure:"length"`
	Reset    bool   `mapstructure:"reset"`
}

type Token struct {
	Type     string `mapstructure:"type"`
	Key      string `mapstructure:"key"`
	Hashalgo string `mapstructure:"hashalgo"`
	Expire   int    `mapstructure:"expire"`
}

type Fastadmin struct {
	Usercenter          bool     `mapstructure:"usercenter"`
	UserRegisterCaptcha string   `mapstructure:"user_register_captcha"`
	LoginCaptcha        bool     `mapstructure:"login_captcha"`
	LoginFailureRetry   bool     `mapstructure:"login_failure_retry"`
	LoginUnique         bool     `mapstructure:"login_unique"`
	LoginipCheck        bool     `mapstructure:"loginip_check"`
	LoginBackground     string   `mapstructure:"login_background"`
	Multiplenav         bool     `mapstructure:"multiplenav"`
	Multipletab         bool     `mapstructure:"multipletab"`
	ShowSubmenu         bool     `mapstructure:"show_submenu"`
	Adminskin           string   `mapstructure:"adminskin"`
	Breadcrumb          bool     `mapstructure:"breadcrumb"`
	Unknownsources      bool     `mapstructure:"unknownsources"`
	BackupGlobalFiles   bool     `mapstructure:"backup_global_files"`
	AutoRecordLog       bool     `mapstructure:"auto_record_log"`
	AddonPureMode       bool     `mapstructure:"addon_pure_mode"`
	CorsRequestDomain   []string `mapstructure:"cors_request_domain"`
	Version             string   `mapstructure:"version"`
	ApiUrl              string   `mapstructure:"api_url"`
}
