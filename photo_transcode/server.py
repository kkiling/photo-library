from flask import Flask, request, jsonify
import numpy as np
from keras.applications import vgg16
from keras.preprocessing.image import load_img, img_to_array
from io import BytesIO
from gevent.pywsgi import WSGIServer

app = Flask(__name__)

preprocess = vgg16.preprocess_input
model = vgg16.VGG16(include_top=False, weights='imagenet', input_tensor=None, input_shape=None, pooling='max')


def get_vector_from_bytes(img_bytes):
    image = load_img(BytesIO(img_bytes), target_size=(224, 224))
    image = img_to_array(image)
    t_arr = np.expand_dims(image, axis=0)
    processed_img = preprocess(t_arr)
    # Берем первый элемент из двухмерного массива
    result = model.predict(processed_img, verbose=0)
    return result[0]


@app.route('/get_vector_from_bytes', methods=['POST'])
def get_vector_from_bytes_endpoint():
    img_bytes = request.data
    vector = get_vector_from_bytes(img_bytes)
    return jsonify(vector.tolist())


if __name__ == '__main__':
    app.run(threaded=True)

