<meta charset="utf-8">
<title>{{.title | default "" | htmlentities}} – {{.site.name|htmlentities}}</title>
<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no">
<meta name="renderer" content="webkit">

{{if .keywords}}
<meta name="keywords" content="{{.keywords|htmlentities}}">
{{end}}
{{if .description}}
<meta name="description" content="{{.description|htmlentities}}">
{{end}}

<link rel="shortcut icon" href="__CDN__/assets/img/favicon.ico" />

<link href="__CDN__/assets/css/frontend{{if eq .Think.config.app_debug false}}.min{{end}}.css?v={{.Think.config.site.version|htmlentities}}" rel="stylesheet">

<!-- HTML5 shim, for IE6-8 support of HTML5 elements. All other JS at the end of file. -->
<!--[if lt IE 9]>
  <script src="__CDN__/assets/js/html5shiv.js"></script>
  <script src="__CDN__/assets/js/respond.min.js"></script>
<![endif]-->
<script type="text/javascript">
    var require = {
        config: {{.Think.config}}
    };
</script>
