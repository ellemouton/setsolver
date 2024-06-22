import cv2
import numpy as np
import sys
import argparse
from PIL import Image
import imutils

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

def draw_contours_on_image(image, contours):
    contour_image = image.copy()

    for c in contours:
        cv2.drawContours(contour_image, [c], 0, (36,255,12), 3)
    
    return contour_image


def edge_detection_from_path(image_path):
    # image is the original image.
    image = read_image(image_path)

    return edge_detection(image)

def edge_detection(image):
    # image_f is the image with filters applied so as to make the edges stand out.
    image_f = apply_filters(image.copy())

    # Find card contours on the image:
    card_contours = find_contours(image_f)
    
    # Create a new image that has the contours drawn on it.
    contours_image = draw_contours_on_image(image, card_contours)
 
    return contours_image

def main():
    parser = argparse.ArgumentParser(description='Edge Detection')
    parser.add_argument('image_path', type=str, help='Path to the input image')
    args = parser.parse_args()

    edge_detection_from_path(args.image_path)

if __name__ == '__main__':
    main()