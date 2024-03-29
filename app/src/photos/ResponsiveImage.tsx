import React, { useState, useEffect, useRef } from 'react';
import { Box } from '@mui/material';
import { PhotoGroupData } from './model';

const imageStyle = {
  width: '100%',
  height: '100%',
  objectFit: 'contain',
};

const ResponsiveImage = ({ photoGroup }: { photoGroup: PhotoGroupData }) => {
  const [load, setLoad] = useState<boolean>(true);
  const ref = useRef<HTMLDivElement>(null);
  const [width, setWidth] = useState<number>(0);

  useEffect(() => {
    const resizeObserver = new ResizeObserver(entries => {
      for (const entry of entries) {
        setWidth(entry.contentRect.width);
      }
    });

    if (ref.current) {
      resizeObserver.observe(ref.current);
    }

    // Отписка от наблюдателя при размонтировании компонента
    return () => {
      if (ref.current) {
        resizeObserver.unobserve(ref.current);
      }
    };
  }, [ref]);

  let imageUrl = photoGroup.original.src;
  for (const preview of photoGroup.previews) {
    if (preview.size >= width) {
      imageUrl = preview.src;
      break;
    }
  }

  return (
    <>
      {load && <div>Загрузка...</div>}
      <Box
        ref={ref}
        component='img'
        src={imageUrl}
        onLoad={() => setLoad(false)}
        sx={{ ...imageStyle, display: load ? 'none' : 'block' }}
      />
    </>
  );
};

export default ResponsiveImage;
