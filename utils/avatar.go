package utils

// List of avatar url and index
var avatarList = []string{
	"https://res.cloudinary.com/dcnuiaskr/image/upload/v1688934868/avatar/avatar1_wklkmm.png",
	"https://res.cloudinary.com/dcnuiaskr/image/upload/v1688934868/avatar/avatar2_lb0x2p.png",
	"https://res.cloudinary.com/dcnuiaskr/image/upload/v1688934868/avatar/avatar3_ah94hx.png",
	"https://res.cloudinary.com/dcnuiaskr/image/upload/v1688934869/avatar/avatar4_pakjxe.png",
	"https://res.cloudinary.com/dcnuiaskr/image/upload/v1688934868/avatar/avatar5_jouzzh.png",
	"https://res.cloudinary.com/dcnuiaskr/image/upload/v1688934868/avatar/avatar6_eg46sj.png",
	"https://res.cloudinary.com/dcnuiaskr/image/upload/v1688934869/avatar/avatar7_gyzakl.png",
	"https://res.cloudinary.com/dcnuiaskr/image/upload/v1688934869/avatar/avatar8_zefxcq.png",
	"https://res.cloudinary.com/dcnuiaskr/image/upload/v1688934869/avatar/avatar9_av5qs9.png",
}

func GetAvatarUrl(index int) string {
	return avatarList[index]
}

