import numpy as np
from keras.applications import vgg16
from keras.preprocessing.image import load_img, img_to_array

# функция создания пакета изображений
# в этом примере обрабатываем по одному изображению
preprocess = vgg16.preprocess_input

# Загрузка модели VGG16 без верхних слоев (полносвязных)
model = vgg16.VGG16(include_top=False, weights='imagenet',
                    input_tensor=None, input_shape=None, pooling='max')


def get_vector(image_path):
    image = load_img(image_path, target_size=(256, 256))
    image = img_to_array(image)
    t_arr = np.expand_dims(image, axis=0)

    processed_img = preprocess(t_arr)

    return model.predict(processed_img, verbose=0)


def similarity(vector1, vector2):
    d1 = np.linalg.norm(vector1, axis=1, keepdims=True)
    d2 = np.linalg.norm(vector2.T, axis=0, keepdims=True)
    return np.dot(vector1, vector2.T) / np.dot(d1, d2)[0][0]

img_path_1 = 'Y:\\documents\\photos_storage\\a0295d5c-09c7-4238-838b-c98e69f9c7bd.jpeg'
img_path_2 = 'Y:\\documents\\photos_storage\\2dc856ba-fe96-49a6-90fd-1e7c9d44acda.jpeg'
#img_path_1 = "Y:\photos\Фото\Архив\Юля древний ноут\Москва осень 2014\DSC_9093.JPG"
#img_path_2 = "Y:\photos\Фото\Архив\Юля древний ноут\Москва осень 2014\DSC_9094.JPG"

v1 = get_vector(img_path_1)
v2 = get_vector(img_path_2)

kef = similarity(v1, v2)

kef = float(kef)
print(kef)
