from flask import Flask, request, send_file, jsonify
from flask_cors import CORS
import cv2
import numpy as np
from PIL import Image, ExifTags
import io
from edge_detect import solve, fill_interpreter

app = Flask(__name__)
CORS(app)  # Enable CORS for all routes

def correct_image_orientation(image):
    try:
        for orientation in ExifTags.TAGS.keys():
            if ExifTags.TAGS[orientation] == 'Orientation':
                break

        exif = dict(image._getexif().items())
        
        if exif[orientation] == 3:
            image = image.rotate(180, expand=True)
        elif exif[orientation] == 6:
            image = image.rotate(270, expand=True)
        elif exif[orientation] == 8:
            image = image.rotate(90, expand=True)
    except (AttributeError, KeyError, IndexError):
        # cases: image don't have getexif
        pass
    return image

@app.route('/health', methods=['GET'])
def health_check():
    thing = fill_interpreter.get_signature_list()

    return jsonify({"status": f'Server is running: {thing}'}), 200

@app.route('/process_image', methods=['POST'])
def process_image():
    if 'image' not in request.files:
        return "No image uploaded", 400

    file = request.files['image']
    image = Image.open(file.stream)
    image = np.array(correct_image_orientation(image))

    print(image.size)
    
    processed_image = solve(image)

    # Save the processed image to a file for debugging
    processed_image_path = 'image.png'
    processed_image_pil = Image.fromarray(processed_image)
    processed_image_pil.save(processed_image_path)

    # Convert the processed image back to a file-like object
    processed_image_pil = Image.fromarray(processed_image)
    buf = io.BytesIO()
    processed_image_pil.save(buf, format='PNG')
    buf.seek(0)

    return send_file(buf, mimetype='image/png')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)