package common

var RootAnnotation = "pvscanner.pvdf.broadsoftware.com"

var PvFreeAnnotation = RootAnnotation + "/free_mib"
var PvSizeAnnotation = RootAnnotation + "/size_mib"

//var NodeVgListAnnotation = RootAnnotation + "/vg_list"
//// %s is volumeGroup name placeholder
//var NodeVgFreeAnnotation = RootAnnotation + "/vg_%s_free"
//var NodeVgSizeAnnotation = RootAnnotation + "/vg_%s_size"

// %s will be replaced by deviceClass
var SizeTopolvmAnnotation = "size.topolvm." + RootAnnotation + "/%s"

