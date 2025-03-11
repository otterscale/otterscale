package model

type Container struct {
	Image           string `json:"image"`
	ImagePullPolicy string `json:"imagePullPolicy"`
}
