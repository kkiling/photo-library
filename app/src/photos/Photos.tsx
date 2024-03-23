import React, { useState, useEffect, useRef } from 'react';
import { Gallery, ThumbnailImageProps } from "react-grid-gallery";
import Lightbox from "yet-another-react-lightbox";
import Counter from "yet-another-react-lightbox/plugins/counter";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import Fullscreen from "yet-another-react-lightbox/plugins/fullscreen";
import "yet-another-react-lightbox/styles.css";
import "yet-another-react-lightbox/plugins/counter.css";
import "yet-another-react-lightbox/plugins/thumbnails.css";
import "yet-another-react-lightbox/styles.css";

const rowHeight = 240

interface PhotoPreview {
  src: string;
  width: number;
  height: number;
}

interface PhotoItem {
  isSelected: boolean;
  src: string;
  width: number;
  height: number;
  previews: PhotoPreview[];
}

const Photos2 = () => {
  const [photos, setPhotos] = useState<PhotoItem[]>([]);
  const [index, setIndex] = useState(-1);
  const [currentPage, setCurrentPage] = useState(1);
  const prevPageRef = useRef(-1);
  const observer = useRef<IntersectionObserver | null>(null);
  const loadingRef = useRef<HTMLDivElement | null>(null);

  const ImageComponent = (props: ThumbnailImageProps) => {
    const { src, alt, style, title } = props.imageProps;
    const {viewportWidth, scaledHeight } = props.item;
    // const item = props.item as any as PhotoItem;
    const size = viewportWidth > scaledHeight ? viewportWidth : scaledHeight;
    const souse = `${src}?size=${size}`
    return <img alt={alt} src={souse} title={title || ""} style={style}  />;
  };

  const slides = photos.map((photo) => ({
    src: photo.src,
    alt: "image 1",
    width: photo.width,
    height: photo.height,
    srcSet: photo.previews.map((preview) => ({
      src: preview.src,
      width: preview.width,
      height: preview.height,
    }))
  }));

  // Загрузка фотографий
  const fetchPhotos = async (page: number) => {
    const response = await fetch(`v1/photos/groups?page=${page}&per_page=24`, {
      method: 'GET',
      headers: {
        'accept': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Ошибка загрузки фотографий');
    }

    const data = await response.json();
    let items: PhotoItem[] = []
    items = data.items.map(( photo: any, index: number) => {
      const main = photo.main_photo;
      const res = {
        isSelected: false,
        src: main.src,
        width: main.width,
        height: main.height,
        previews: main.previews.map((preview: any) => ({
          src: preview.src,
          width: preview.width,
          height: preview.height,
        }))
      };
      return res;
    });

    return items
  };

  // Подгрузка страницы
  useEffect(() => {
    const loadPhotos = async () => {
      const newPhotos = await fetchPhotos(currentPage);
      setPhotos(prevPhotos => [...prevPhotos, ...newPhotos]);
    };

    console.log(prevPageRef.current)
    if (currentPage > prevPageRef.current) {
      loadPhotos();
      prevPageRef.current = currentPage;
    }
  }, [currentPage]);

  const onSelect = (index: number, item: PhotoItem) => {
    const nextImages = photos.map((image, i) =>
      i === index ? { ...image, isSelected: !image.isSelected } : image
    );
    setPhotos(nextImages);
  };

  const onClick = (index: number, item: PhotoItem) => setIndex(index);

  // Ленивая подгрузка в момент скрола к концу страницы
  useEffect(() => {
    const options = {
      root: null,
      rootMargin: '0px',
      threshold: 1.0
    };

    observer.current = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting) {
        console.log("new page")
        setCurrentPage((prevPage) => prevPage + 1);
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
  }, []);

  return (
    <div id="photos-container">
      <Gallery id="gallery" images={photos} 
        rowHeight={rowHeight} 
        onSelect={onSelect}
        onClick={onClick}
        enableImageSelection={true}
        thumbnailImageComponent={ImageComponent}
      />
      <Lightbox
        slides={slides}
        open={index >= 0}
        index={index}
        close={() => setIndex(-1)}
        on={{
          view: (index) => console.log("View", index),
          entering: () => console.log("Entering"),
          entered: () => console.log("Entered"),
          exiting: () => console.log("Exiting"),
          exited: () => console.log("Exited"),
        }}
        plugins={[Zoom, Fullscreen, Thumbnails, Counter]}
        counter={{ container: { style: { top: "unset", bottom: 0 } } }}
        thumbnails={{ position: "bottom"} }
      />
      <div ref={loadingRef} style={{ height: "100px", margin: "0 auto" }}>
        Загрузка...
      </div>
    </div>
  );
};

export default Photos2;
