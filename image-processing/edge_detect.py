import cv2
import numpy as np
import sys
import argparse
from PIL import Image
import imutils
import os
from card import Card
import tensorflow as tf
import requests

url = "http://localhost:8081/set"

TF_FILL_MODEL_FILE_PATH = 'fill_model.tflite' 

fill_interpreter = tf.lite.Interpreter(model_path=TF_FILL_MODEL_FILE_PATH)
fill_classify_lite = fill_interpreter.get_signature_runner('serving_default')

fill_input_details = fill_interpreter.get_input_details()
fill_output_details = fill_interpreter.get_output_details()

TF_SHAPE_MODEL_FILE_PATH = 'shape_model.tflite' 

shape_interpreter = tf.lite.Interpreter(model_path=TF_SHAPE_MODEL_FILE_PATH)
shape_classify_lite = shape_interpreter.get_signature_runner('serving_default')

shape_input_details = shape_interpreter.get_input_details()
shape_output_details = shape_interpreter.get_output_details()

TF_COLOUR_MODEL_FILE_PATH = 'colour_model.tflite' 

colour_interpreter = tf.lite.Interpreter(model_path=TF_COLOUR_MODEL_FILE_PATH)
colour_classify_lite = colour_interpreter.get_signature_runner('serving_default')

colour_input_details = colour_interpreter.get_input_details()
colour_output_details = colour_interpreter.get_output_details()

def read_image(image_path):
    # Load the image using OpenCV.
    image = cv2.imread(image_path)
    if image is None:
        print(f"Error: Unable to load image {image_path}")
        return

    return image

def apply_filters(image):
    # Apply grey scale filter.
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

    # Gaussian blur the image.
    blur = cv2.GaussianBlur(gray, (3, 3), 0)

    # Otsu's threshold
    thresh = cv2.threshold(blur, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)[1]

    return thresh

def rotate_and_crop(image, rect):
    """
    Rotate the image around the center of the rectangle and crop it,
    ensuring the longer side is on the x-axis.

    Parameters:
    - image: The source image.
    - rect: The bounding rectangle (center, size, angle).

    Returns:
    - The rotated and cropped image.
    """
    center, size, angle = rect[0], rect[1], rect[2]
    width, height = size

    # Ensure the longer side is on the x-axis
    if width < height:
        angle += 90
        size = (height, width)

    size = tuple(map(int, size))

    # Get rotation matrix
    M = cv2.getRotationMatrix2D(center, angle, 1.0)
    height, width = image.shape[:2]

    # Perform rotation
    rotated = cv2.warpAffine(image, M, (width, height))

    # Crop the rotated rectangle
    cropped = cv2.getRectSubPix(rotated, size, center)

    return cropped

# We don't want to capture the smaller shapes, only the cards. So set a 
# threshold here. 
# TODO: make this dynamic (ie, find contours and then choose threshold 
#       to filter out the cards)
threshold_min_area = 9000

def is_card_contour(contour):
    area = cv2.contourArea(contour)

    return area > threshold_min_area

def find_contours(image):
    # Find contours and filter for cards using contour area.
    # RETR_EXTERNAL: only tries to find extreme outer contours.
    cnts = cv2.findContours(image, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    # Get the contours list. Depending on OpenCV version, this is a different 
    # in the return list.
    cnts = cnts[0] if len(cnts) == 2 else cnts[1]

    # Filter out card contours.
    card_contours = [cnt for cnt in cnts if is_card_contour(cnt)]

    print("found ", len(card_contours), "cards")

    return card_contours

def draw_cards_on_image(image, cards, setIDs):
    contour_image = image.copy()

    for card in cards:
        if card.id in setIDs:
            cv2.drawContours(contour_image, [card.contour], 0, (36,255,12), 3)
        
        card.writeOnImage(contour_image)
    
    return contour_image

def writeOnImage(image, text):
    cv2.putText(image, text, (100,100), cv2.FONT_HERSHEY_SIMPLEX, 2, (0, 0, 0), 2, cv2.LINE_AA)

def edge_detection_from_path(image_path):
    # image is the original image.
    image = read_image(image_path)

    return find_and_draw_cards(image)

def card_to_dict(card):
    return {
        "id": card.id,
        "colour": card.colour,
        "fill": card.fill,
        "shape": card.shape,
        "count": card.count
    }

def solve(image):
    # Find the cards on the image and classify them all.
    cards = find_cards(image)

    card_map = {card.id: card for card in cards}
    cards_dict = [card_to_dict(card) for card in cards]

    response = requests.post(url, json={"cards": cards_dict})

    if response.status_code != 200:
        print(f"Failed to get response. Status code: {response.status_code}")

        writeOnImage(image, "failed to get response from solver")

        return image

    data = response.json()
    sets = data['sets']

    if len(sets) == 0:
        writeOnImage(image, "no SET found!")

        return image

 
    #for i, set_ in enumerate(sets):
    #    print(f"Set {i+1}:")

    # Create a new image that has the contours drawn on it.
    contours_image = draw_cards_on_image(image, cards, sets[0])

    return contours_image


def find_and_draw_cards(image):
    cards = find_cards(image)
    
    # Create a new image that has the contours drawn on it.
    contours_image = draw_cards_on_image(image, cards)
 
    return contours_image

def find_cards(image):
    # image_f is the image with filters applied so as to make the edges stand out.
    image_f = apply_filters(image.copy())

    # Find card contours on the image:
    card_contours = find_contours(image_f)
    
    # For each card contour, create a new Card object.
    cards = []
    for contour in card_contours: 
        card = Card(contour, image)

        # First find the count.
        cardCount(card)

        # Attempt to guess the fill.
        cardFill(card)
        
        cardShape(card)

        cardColour(card)

        cards.append(card)

    return cards

colour_class_names = ['green', 'purple', 'red']

def cardColour(card):
    img_height, img_width = colour_input_details[0]['shape'][1], colour_input_details[0]['shape'][2]
    single_shape = cv2.resize(card.single_shape, (img_height, img_width))

    img_array = tf.keras.utils.img_to_array(single_shape)
    img_array = tf.expand_dims(img_array, 0) # Create a batch

    predictions_lite = colour_classify_lite(sequential_6_input=img_array)['outputs']
    score_lite = tf.nn.softmax(predictions_lite)

    colour = colour_class_names[np.argmax(score_lite)]

    card.setColour(colour)

shape_class_names = ['diamond', 'oval', 'squiggle']

def cardShape(card):
    img_height, img_width = shape_input_details[0]['shape'][1], shape_input_details[0]['shape'][2]
    single_shape = cv2.resize(card.single_shape, (img_height, img_width))

    img_array = tf.keras.utils.img_to_array(single_shape)
    img_array = tf.expand_dims(img_array, 0) # Create a batch

    predictions_lite = shape_classify_lite(sequential_5_input=img_array)['outputs']
    score_lite = tf.nn.softmax(predictions_lite)

    shape = shape_class_names[np.argmax(score_lite)]

    card.setShape(shape)

fill_class_names = ['hollow', 'shaded', 'solid']

def cardFill(card):
    img_height, img_width = fill_input_details[0]['shape'][1], fill_input_details[0]['shape'][2]
    single_shape = cv2.resize(card.single_shape, (img_height, img_width))

    img_array = tf.keras.utils.img_to_array(single_shape)
    img_array = tf.expand_dims(img_array, 0) # Create a batch

    predictions_lite = fill_classify_lite(sequential_10_input=img_array)['outputs']
    score_lite = tf.nn.softmax(predictions_lite)

    fill = fill_class_names[np.argmax(score_lite)]

    card.setFill(fill)
    

def cardCount(card):
    # Some filtering to make the shape stand out.
    gray = cv2.cvtColor(card.image, cv2.COLOR_BGR2GRAY)
    blur = cv2.GaussianBlur(gray, (3, 3), 0)
    thresh = cv2.threshold(blur, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)[1]

    # Find contours and filter for cards using contour area.
    # RETR_TREE: tries to find all contours.
    cnts, hierarchy = cv2.findContours(thresh, cv2.RETR_TREE, cv2.CHAIN_APPROX_SIMPLE)
    contours_image = card.image.copy()

    # Filter out smaller and larger contours.
    threshold_min_area = 400
    threshold_max_area = 15000
    finalContours = []
    for i, contour in enumerate(cnts):
        parent_index = hierarchy[0][i][3]
        hasParent = parent_index != -1

        area = cv2.contourArea(contour)

        # Filter out an area that is probably a card. 
        if area > threshold_max_area:
            continue

        cv2.drawContours(contours_image, [contour], 0, (36,255,12), 3)

        # Area of this contour itself is too small/large.
        if area < threshold_min_area: 
            continue

        # Usually we dont want to include contours that have a parent
        # unless that parent is the card itself.
        if hasParent and cv2.contourArea(cnts[parent_index]) < threshold_max_area:
            continue

        finalContours.append(contour)

    count = len(finalContours)
    if count == 0:
        raise Exception("No shapes were found on the card!")

    card.setCount(count)

    # Now we can just focus on one shape. So select one contour
    # and zoom in on that. 
    contour = finalContours[0] 

    card.zoomOnShape(contour)


pic_counter = 0
def create_training_set(image):
    global pic_counter
    pic_counter+=1

    cards = find_cards(image)

    output_dir = 'training_sets'
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    for i, card in enumerate(cards):
        # Save the image to a file
        output_path = os.path.join(output_dir, f'card_{pic_counter}_{i+1}.png')
        cv2.imwrite(output_path, card.image)
    
def main():
    parser = argparse.ArgumentParser(description='Edge Detection')
    parser.add_argument('image_path', type=str, help='Path to the input image')
    args = parser.parse_args()

    edge_detection_from_path(args.image_path)

if __name__ == '__main__':
    main()