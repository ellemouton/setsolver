import cv2
from PIL import Image
import os

pic_counter = 0

class Card:
    def __init__(self, contour, image):
        # Contour is the original contour of the card on the image.
        self.contour = contour

        # rectangle is the estimated rectangle representing
        # the card on the image.
        self.rect = cv2.minAreaRect(contour)

        # image is an image of just this card. This will 
        # be used during classification of the card.
        self.image = rotate_and_crop(image, self.rect)

    def setColour(self, colour): 
        self.colour = colour

    def setShape(self, shape): 
        self.shape = shape

    def setFill(self, fill):
        self.fill = fill
    
    def setCount(self, count):
        self.count = count

    def writeOnImage(self, image):
        # Calculate the center of the rectangle
        center_x = int(self.rect[0][0])
        center_y = int(self.rect[0][1])

        count_offset_x = -60
        count_offset_y = -35

        colour_offset_x = -60
        colour_offset_y = -5

        fill_offset_x = -60
        fill_offset_y = 25
        
        shape_offset_x = -60
        shape_offset_y = 50
        
        # Calculate the position to put the text (centered)
        count_position = (center_x + count_offset_x, center_y + count_offset_y)  
        colour_position = (center_x + colour_offset_x, center_y + colour_offset_y)  
        fill_position = (center_x + fill_offset_x, center_y + fill_offset_y)  
        shape_position = (center_x + shape_offset_x, center_y + shape_offset_y)  
        
        # Write count.
        cv2.putText(image, f'{self.count}', count_position, cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 0), 2, cv2.LINE_AA)
        cv2.putText(image, f'{self.colour}', colour_position, cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 0), 2, cv2.LINE_AA)
        cv2.putText(image, f'{self.fill}', fill_position, cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 0), 2, cv2.LINE_AA)
        cv2.putText(image, f'{self.shape}', shape_position, cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 0), 2, cv2.LINE_AA)

    def __str__(self):
        return f"{self.count}, {self.colour}, {self.fill}, {self.shape}"

    def zoomOnShape(self, contour):
        global pic_counter
        pic_counter+=1

        x, y, w, h = cv2.boundingRect(contour)
        self.single_shape = self.image[y:y+h, x:x+w]

        output_dir = 'temp'
        if not os.path.exists(output_dir):
           os.makedirs(output_dir)

        # Save the processed image to a file for debugging
        output_path = os.path.join(output_dir, f'shape_{pic_counter}.png')
        image_to_save = cv2.cvtColor(self.single_shape, cv2.COLOR_BGR2RGB)
        cv2.imwrite(output_path, image_to_save)

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