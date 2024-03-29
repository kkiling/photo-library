export interface PhotoPreview {
  src: string;
  width: number;
  height: number;
  size: number;
}

export interface Tag {
  id: string;
  name: string;
  type: string;
  color: string;
}

export interface Geo {
  latitude: number;
  longitude: number;
}

export interface MetaData {
  modelInfo: string;
  sizeBytes: number;
  widthPixel: number;
  heightPixel: number;
  dataTime: string;
  updateAt: string;
  geo: Geo | null;
}

export interface PhotoGroupData {
  groupId: string;
  photosCount: number;
  original: PhotoPreview;
  previews: PhotoPreview[];
  tags: Tag[];
  metaData: MetaData;
}

export interface PhotoGalleryData {
  isSelected: boolean;
  src: string;
  width: number;
  height: number;

  groupId: string;
  photosCount: number;
  original: PhotoPreview;
  previews: PhotoPreview[];
}