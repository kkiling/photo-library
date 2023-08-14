-- +goose Up
-- +goose StatementBegin
CREATE TABLE exif_data (
   photo_id UUID PRIMARY KEY,
   sharpness INT,
   maker_note VARCHAR(256),
   x_resolution FLOAT8,
   exposure_bias_value FLOAT8,
   aperture_value FLOAT8,
   resolution_unit INT,
   thumb_jpeg_interchange_format INT,
   model VARCHAR(256),
   focal_plane_x_resolution FLOAT8,
   bits_per_sample INT[],
   device_setting_description VARCHAR(256),
   sub_sec_time_digitized VARCHAR(256),
   image_length INT,
   gps_version_id INT,
   xp_author INT[],
   max_aperture_value FLOAT8,
   iso_speed_ratings INT,
   sensing_method INT,
   file_source VARCHAR(256),
   exif_version VARCHAR(256),
   image_width INT,
   gps_altitude_ref FLOAT8,
   samples_per_pixel INT,
   software VARCHAR(256),
   saturation INT,
   light_source INT,
   gps_latitude FLOAT8[],
   lens_model VARCHAR(256),
   pixel_x_dimension INT,
   user_comment VARCHAR(256),
   gps_info_ifd_pointer INT,
   exif_ifd_pointer INT,
   exposure_time FLOAT8,
   cfa_pattern VARCHAR(256),
   gps_longitude FLOAT8[],
   copyright VARCHAR(256),
   flash INT,
   y_resolution FLOAT8,
   metering_mode INT,
   interoperability_ifd_pointer INT,
   gps_altitude FLOAT8,
   scene_type VARCHAR(256),
   image_description VARCHAR(256),
   brightness_value FLOAT8,
   artist VARCHAR(256),
   focal_length_in_35_mm_film INT,
   exposure_mode INT,
   white_balance INT,
   pixel_y_dimension INT,
   sub_sec_time_original VARCHAR(256),
   digital_zoom_ratio FLOAT8,
   compressed_bits_per_pixel FLOAT8,
   date_time VARCHAR(256),
   gain_control INT,
   xp_keywords INT[],
   contrast INT,
   gps_processing_method VARCHAR(256),
   focal_plane_resolution_unit INT,
   shutter_speed_value FLOAT8,
   sub_sec_time VARCHAR(256),
   gps_latitude_ref VARCHAR(256),
   image_unique_id VARCHAR(256),
   color_space INT,
   exposure_program INT,
   gps_longitude_ref VARCHAR(256),
   date_time_original VARCHAR(256),
   components_configuration VARCHAR(256),
   f_number FLOAT8,
   xp_comment INT[],
   xp_title INT[],
   thumb_jpeg_interchange_format_length INT,
   y_cb_cr_positioning INT,
   orientation INT,
   photometric_interpretation INT,
   focal_length FLOAT8,
   scene_capture_type INT,
   make VARCHAR(256),
   gps_date_stamp VARCHAR(256),
   related_sound_file VARCHAR(256),
   flashpix_version VARCHAR(256),
   custom_rendered INT,
   date_time_digitized VARCHAR(256),
   exposure_index FLOAT8,
   focal_plane_y_resolution FLOAT8,
   interoperability_index VARCHAR(256),
   subject_distance_range INT,
   gps_time_stamp FLOAT8[]
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE exif_data;
-- +goose StatementEnd