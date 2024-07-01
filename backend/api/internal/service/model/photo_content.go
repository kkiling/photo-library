package model

// PhotoContentDTO контент фото и ее расширение
type PhotoContentDTO struct {
	PhotoBody []byte
	Extension PhotoExtension
}
