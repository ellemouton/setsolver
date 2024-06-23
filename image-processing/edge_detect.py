import cv2
import numpy as np
import sys
import argparse
from PIL import Image
import imutils
import os
from card import Card

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

def draw_cards_on_image(image, cards):
    contour_image = image.copy()

    for card in cards:
        cv2.drawContours(contour_image, [card.contour], 0, (36,255,12), 3)
        card.writeOnImage(contour_image)
    
    return contour_image


def edge_detection_from_path(image_path):
    # image is the original image.
    image = read_image(image_path)

    return find_and_draw_cards(image)

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
    j = 0
    for contour in card_contours: 
        card = Card(contour, image)

        # First find the count.
        cardCount(card, j)
        j+=1

        cards.append(card)

    return cards


def cardCount(card, j):
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