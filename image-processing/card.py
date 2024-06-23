import cv2

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
    
    def setCount(self, count):
        self.count = count

    def writeOnImage(self, image):
        # Calculate the center of the rectangle
        center_x = int(self.rect[0][0])
        center_y = int(self.rect[0][1])
        
        # Calculate the position to put the text (centered)
        text_position = (center_x - 10, center_y + 10)  # Adjust offsets as necessary
        
        # Put the text in the rectangle
        cv2.putText(image, f"{self.count}", text_position, cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 0), 2, cv2.LINE_AA)


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