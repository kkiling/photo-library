import React, { useState, useEffect, useRef } from 'react';
import { Gallery, ThumbnailImageProps } from 'react-grid-gallery';
import PhotoGroupView from './PhotoGroupView';
import { PhotoGalleryData } from './model';

const rowHeight = 240;
const perPage = 32;

const fetchGroups = async (page: number) => {
  const response = await fetch(`v1/photos/groups?page=${page}&per_page=${perPage}`, {
    method: 'GET',
    headers: {
      accept: 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error('Ошибка загрузки фотографий');
  }

  const data = await response.json();
  let groups: PhotoGalleryData[] = [];

  groups = data.items.map((group: any) => {
    const res = {
      isSelected: false,
      src: group.original.src,
      width: group.original.width,
      height: group.original.height,
      groupId: group.id,
      photoCount: group.photo_count,
      original: {
        src: group.original.src,
        width: group.original.width,
        height: group.original.height,
        size: group.original.size,
      },
      previews: group.previews.map((preview: any) => ({
        src: preview.src,
        width: preview.width,
        height: preview.height,
        size: preview.size,
      })),
    };
    return res;
  });

  return groups;
};


const ImageComponent = (props: ThumbnailImageProps) => {
  const { alt, style, title } = props.imageProps;
  const { viewportWidth, scaledHeight } = props.item;

  const item = props.item as unknown as PhotoGalleryData;
  const viewportSize =
    viewportWidth > scaledHeight ? viewportWidth : scaledHeight;

  let imageUrl = item.original.src;
  for (const preview of item.previews) {
    if (preview.size >= viewportSize) {
      imageUrl = preview.src;
      break;
    }
  }

  return (
    <img alt={alt} src={imageUrl} title={title || ''} style={{ ...style }} />
  );
};

const PhotoGallery = () => {
  const [load, setLoad] = useState<boolean>(true);
  const [groups, setGroups] = useState<PhotoGalleryData[]>([]);
  const [index, setIndex] = useState<number | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const prevPageRef = useRef(-1);
  const observer = useRef<IntersectionObserver | null>(null);
  const loadingRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const loadPhotos = async () => {
      setLoad(true);
      const groups = await fetchGroups(currentPage);
      setLoad(false);
      setGroups(prevGroups => [...prevGroups, ...groups]);
      if (index != null) {
        setIndex(index + 1);
      }
    };

    if (currentPage > prevPageRef.current) {
      loadPhotos();
      prevPageRef.current = currentPage;
    }
  }, [currentPage]);

  const onClick = (index: number) => setIndex(index);

  useEffect(() => {
    const options = {
      root: null,
      rootMargin: '0px',
      threshold: 1.0,
    };

    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting) {
        setCurrentPage(prevPage => prevPage + 1);
      }
    }, options);

    if (loadingRef.current) {
      observer.current.observe(loadingRef.current);
    }

    return () => {
      if (observer.current && loadingRef.current) {
        observer.current.unobserve(loadingRef.current);
      }
    };
    
  }, [groups]);

  const onNext = () => {
    if (index == null) return;
    if (index == groups.length - 1) {
      setCurrentPage(prevPage => prevPage + 1);
    } else {
      setIndex(index + 1);
    }
  };

  const onPrev = () => {
    if (index == null || index == 0) return;
    setIndex(index - 1);
  };

  if (groups.length == 0) return <div>Загрузка...</div>;

  return (
    <div id='photos-container'>
      <PhotoGroupView
        parentLoad={load}
        groupId={index == null ? null : groups[index].groupId}
        onClose={() => setIndex(null)}
        onNext={onNext}
        onPrev={onPrev}
      />
      <Gallery
        id='gallery'
        images={groups}
        rowHeight={rowHeight}
        onClick={onClick}
        enableImageSelection={true}
        thumbnailImageComponent={ImageComponent}
      />
      <div ref={loadingRef} style={{ margin: 'auto' }}>
        {load ? 'Загрузка...' : ''}
      </div>
    </div>
  );
};

export default PhotoGallery;
