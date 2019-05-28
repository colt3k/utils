package mimetypes

var types map[string]string

func Find(key string ) string {

	loadTypes()
	return types[key]
}

func loadTypes()  {
	if types == nil {
		types = make(map[string]string)
		types["7z"] = "application/x-7z-compressed"
		types["bz"] = "application/x-bzip"
		types["bz2"] = "application/x-bzip2"
		types["zip"] = "application/zip"
		types["exe"] = "application/x-msdownload"
		types["txt"] = "text/plain"
		types["html"] = "text/html"
		types["css"] = "text/css"
		types["csv"] = "text/csv"
		types["ics"] = "text/calendar"
		types["js"] = "application/javascript"
		types["json"] = "application/json"
		types["jpeg"] = "image/jpeg"
		types["jpg"] = "image/jpeg"
		types["xls"] = "application/vnd.ms-excel"
		types["xlsx"] = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		types["doc"] = "application/msword"
		types["docx"] = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		types["ppt"] = "application/vnd.ms-powerpoint"
		types["vsd"] = "application/vnd.visio"
		types["mid"] = "audio/midi"
		types["mpeg"] = "video/mpeg"
		types["mpg"] = "video/mpeg"
		types["mp4"] = "video/mp4"
		types["ogx"] = "application/ogg"
		types["oga"] = "audio/ogg"
		types["ogv"] = "video/ogg"
		types["psd"] = "image/vnd.adobe.photoshop"
		types["pgp"] = "application/pgp-encrypted"
		types["qt"] = "video/quicktime"
		types["rar"] = "application/x-rar-compressed"
		types["rtx"] = "text/richtext"
		types["svg"] = "image/svg+xml"
		types["au"] = "audio/basic"
		types["tiff"] = "image/tiff"
		types["tar"] = "application/x-tar"
		types["wav"] = "audio/x-wav"
		types["wsdl"] = "application/wsdl+xml"
		types["xhtml"] = "application/xhtml+xml"
		types["xml"] = "application/xml"
		types["xslt"] = "application/xslt+xml"
		types["yaml"] = "text/yaml"
	}
}
