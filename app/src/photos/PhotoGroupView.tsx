import React, { useState, useEffect } from 'react';
import { Modal, Box, Typography, IconButton, Button } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import ArrowBackIosNewIcon from '@mui/icons-material/ArrowBackIosNew';
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';
import { PhotoGroupData } from './model';
import ResponsiveImage from './ResponsiveImage';
import PhotoProperty from './PhotoProperty';

const modalStyle = {
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
};

const mainBoxStyle = {
  width: '95%',
  height: '95%',
  display: 'flex',
  bgcolor: 'background.paper',
};

const imageBoxStyle = {
  position: 'relative',
  width: '60%',
  height: '100%',
};

const propertyBoxStyle = {
  width: '40%',
  height: '100%',
  padding: 2,
  display: 'flex',
  flexDirection: 'column',
  borderLeft: 1,
  borderColor: 'divider',
  p: 0,
};

const propertyTopStyle = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  borderBottom: 1,
  borderColor: 'divider',
  p: 1,
};

const propertyBottomStyle = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  borderTop: 1,
  borderColor: 'divider',
  p: 1,
};

const propertyContentStyle = {
  p: 2,
  flexGrow: 1,
};

const fetchPhotos = async (groupId: string) => {
  const response = await fetch(`v1/photos/groups/${groupId}`, {
    method: 'GET',
    headers: {
      accept: 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error('Ошибка загрузки фотографий');
  }

  const data = await response.json();
  const res: PhotoGroupData = {
    groupId: data.id,
    photosCount: data.photo_count,
    original: {
      src: data.original.src,
      width: data.original.width,
      height: data.original.height,
      size: data.original.size,
    },
    previews: data.previews.map((preview: any) => ({
      src: preview.src,
      width: preview.width,
      height: preview.height,
      size: preview.size,
    })),
    tags: data.tags.map((tag: any) => ({
      id: tag.id,
      name: tag.name,
      type: tag.type,
      color: tag.color,
    })),
    metaData: {
      modelInfo: data.meta_data.model_info,
      sizeBytes: data.meta_data.size_bytes,
      widthPixel: data.meta_data.width_pixel,
      heightPixel: data.meta_data.height_pixel,
      dataTime: data.meta_data.data_time,
      updateAt: data.meta_data.update_at,
      geo: data.geo ? {
        latitude: data.geo.latitude,
        longitude: data.geo.longitude,
      } : null,
    }
  };

  return res;
};

interface Property {
  parentLoad: boolean;
  groupId: string | null;
  onClose: () => void;
  onPrev: () => void;
  onNext: () => void;
}

const PhotoGroupView = ({
  parentLoad,
  groupId,
  onClose,
  onPrev,
  onNext,
}: Property) => {
  const [load, setLoad] = useState<boolean>(true);
  const [photoGroup, setPhotoGroup] = useState<PhotoGroupData | null>(null);

  // Подгрузка страницы
  useEffect(() => {
    const loadPhotos = async () => {
      setLoad(true);
      const photoGroup = await fetchPhotos(groupId!);
      setLoad(false);
      setPhotoGroup(photoGroup);
    };
    if (groupId != null) loadPhotos();
  }, [groupId]);

  const isLoad = load || parentLoad;

  return (
    <Modal
      open={groupId != null}
      onClose={onClose}
      aria-labelledby='full-screen-modal-title'
      aria-describedby='full-screen-modal-description'
      sx={modalStyle}
    >
      <Box sx={mainBoxStyle}>
        <Box sx={imageBoxStyle}>
          {isLoad ? (
            <div>Загрузка...</div>
          ) : (
            <ResponsiveImage photoGroup={photoGroup!} />
          )}
        </Box>
        {/* Панель с контролями и текстом */}
        <Box sx={propertyBoxStyle}>
          {/* Верхняя часть */}
          <Box sx={propertyTopStyle}>
            <Typography variant='h6'>Заголовок</Typography>
            <IconButton onClick={onClose} sx={{ alignSelf: 'flex-end' }}>
              <CloseIcon />
            </IconButton>
          </Box>
          {/* Контент */}
          <Box sx={propertyContentStyle}>
            {/* Сюда добавьте ваш контент */}
            {isLoad ? <div>Загрузка...</div> : <PhotoProperty photoGroup={photoGroup!}/>}
          </Box>
          {/* Нижняя часть */}
          <Box sx={propertyBottomStyle}>
            {/* Кнопки влево и вправо */}
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              {/* Добавлено для лучшего выравнивания кнопок */}
              <IconButton
                disabled={isLoad}
                onClick={onPrev}
                sx={{ alignSelf: 'flex-end', mr: 1 }}
              >
                <ArrowBackIosNewIcon />
              </IconButton>
              <IconButton
                disabled={isLoad}
                onClick={onNext}
                sx={{ alignSelf: 'flex-end' }}
              >
                <ArrowForwardIosIcon />
              </IconButton>
            </Box>
            {/* Кнопки отмена и сохранить */}
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              <Button disabled={isLoad} variant='outlined' sx={{ mr: 1 }}>
                Отмена
              </Button>
              <Button disabled={isLoad} variant='contained'>
                Сохранить
              </Button>
            </Box>
          </Box>
        </Box>
      </Box>
    </Modal>
  );
};

export default PhotoGroupView;
