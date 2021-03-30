'''
    Text detection model: https://github.com/argman/EAST
    Download link: https://www.dropbox.com/s/r2ingd0l3zt8hxs/frozen_east_text_detection.tar.gz?dl=1
    CRNN Text recognition model taken from here: https://github.com/meijieru/crnn.pytorch
    How to convert from pb to onnx:
    Using classes from here: https://github.com/meijieru/crnn.pytorch/blob/master/models/crnn.py
    More converted onnx text recognition models can be downloaded directly here:
    Download link: https://drive.google.com/drive/folders/1cTbQ3nuZG-EKWak6emD_s8_hHXWz7lAr?usp=sharing
    And these models taken from here:https://github.com/clovaai/deep-text-recognition-benchmark
    import torch
    from models.crnn import CRNN
    model = CRNN(32, 1, 37, 256)
    model.load_state_dict(torch.load('crnn.pth'))
    dummy_input = torch.randn(1, 1, 32, 100)
    torch.onnx.export(model, dummy_input, "crnn.onnx", verbose=True)
    python test.py --input .\FRAMES\CN1_3.jpg -m .\frozen_east_text_detection.pb
'''

# Import required modules
import numpy as np
import cv2 as cv
import math
import argparse
from sklearn.cluster import DBSCAN
from pprint import pprint
import os

############ Add argument parser for command line arguments ############
parser = argparse.ArgumentParser(
    description="Use this script to run TensorFlow implementation (https://github.com/argman/EAST) of "
                "EAST: An Efficient and Accurate Scene Text Detector (https://arxiv.org/abs/1704.03155v2)"
                "The OCR model can be obtained from converting the pretrained CRNN model to .onnx format from the github repository https://github.com/meijieru/crnn.pytorch"
                "Or you can download trained OCR model directly from https://drive.google.com/drive/folders/1cTbQ3nuZG-EKWak6emD_s8_hHXWz7lAr?usp=sharing")
parser.add_argument('--input',required=True,
                    help='Path to input image.')
parser.add_argument('--output',required=True,
                    help='Output path', default="test_overlay.jpg")
parser.add_argument('--model', '-m', default="frozen_east_text_detection.pb",
                    help='Path to a binary .pb file contains trained detector network.')
parser.add_argument('--width', type=int, default=320,
                    help='Preprocess input image by resizing to a specific width. It should be multiple by 32.')
parser.add_argument('--height', type=int, default=320,
                    help='Preprocess input image by resizing to a specific height. It should be multiple by 32.')
parser.add_argument("--list", "-l", nargs="*", default=[],
                    help="List of points to be higlighted")
parser.add_argument('--generate_reference', "-gr", action='store_true',
                    help="Generate a reference image for highlighting")
parser.add_argument('--thr', type=float, default=0.5,
                    help='Confidence threshold.')
parser.add_argument('--dist', type=int, default=30,
                    help='Box merge threshold distance')
parser.add_argument('--nms', type=float, default=0.4,
                    help='Non-maximum suppression threshold.')
args = parser.parse_args()


############ Utility functions ############


def decodeBoundingBoxes(scores, geometry, scoreThresh):
    detections = []
    confidences = []

    ############ CHECK DIMENSIONS AND SHAPES OF geometry AND scores ############
    assert len(scores.shape) == 4, "Incorrect dimensions of scores"
    assert len(geometry.shape) == 4, "Incorrect dimensions of geometry"
    assert scores.shape[0] == 1, "Invalid dimensions of scores"
    assert geometry.shape[0] == 1, "Invalid dimensions of geometry"
    assert scores.shape[1] == 1, "Invalid dimensions of scores"
    assert geometry.shape[1] == 5, "Invalid dimensions of geometry"
    assert scores.shape[2] == geometry.shape[2], "Invalid dimensions of scores and geometry"
    assert scores.shape[3] == geometry.shape[3], "Invalid dimensions of scores and geometry"
    height = scores.shape[2]
    width = scores.shape[3]
    for y in range(0, height):

        # Extract data from scored
        scoresData = scores[0][0][y]
        x0_data = geometry[0][0][y]
        x1_data = geometry[0][1][y]
        x2_data = geometry[0][2][y]
        x3_data = geometry[0][3][y]
        anglesData = geometry[0][4][y]
        for x in range(0, width):
            score = scoresData[x]

            # If score is lower than threshold score, move to next x
            if (score < scoreThresh):
                continue

            # Calculate offset
            offsetX = x * 4.0
            offsetY = y * 4.0
            angle = anglesData[x]

            # Calculate cos and sin of angle
            cosA = math.cos(angle)
            sinA = math.sin(angle)
            h = x0_data[x] + x2_data[x]
            w = x1_data[x] + x3_data[x]

            # Calculate offset
            offset = ([offsetX + cosA * x1_data[x] + sinA * x2_data[x],
                       offsetY - sinA * x1_data[x] + cosA * x2_data[x]])

            # Find points for rectangle
            p1 = (-sinA * h + offset[0], -cosA * h + offset[1])
            p3 = (-cosA * w + offset[0], sinA * w + offset[1])
            center = (0.5 * (p1[0] + p3[0]), 0.5 * (p1[1] + p3[1]))
            detections.append((center, (w, h), -1 * angle * 180.0 / math.pi))
            confidences.append(float(score))

    # Return detections and confidences
    return [detections, confidences]


def east(detector, frame, inpWidth, inpHeight):
    outNames = []
    outNames.append("feature_fusion/Conv_7/Sigmoid")
    outNames.append("feature_fusion/concat_3")

    detector.setInput(frame)
    return detector.forward(outNames)


def resize_image_blob(im, max_side_len=2400):
    '''
    resize image to a size multiple of 32 which is required by the network
    :param im: the resized image
    :param max_side_len: limit of max image size to avoid out of memory in gpu
    :return: the resized image and the resize ratio
    '''
    h, w, _ = im.shape
    resize_w = w
    resize_h = h

    # limit the max side
    if max(resize_h, resize_w) > max_side_len:
        ratio = float(
            max_side_len) / resize_h if resize_h > resize_w else float(max_side_len) / resize_w
    else:
        ratio = 1.
    resize_h = int(resize_h * ratio)
    resize_w = int(resize_w * ratio)

    resize_h = resize_h if resize_h % 32 == 0 else (resize_h // 32) * 32
    resize_w = resize_w if resize_w % 32 == 0 else (resize_w // 32) * 32
    resize_h = max(32, resize_h)
    resize_w = max(32, resize_w)
    # im = cv.resize(im, (int(resize_w), int(resize_h)))
    blob = cv.dnn.blobFromImage(
        im, 1.0, (int(resize_w), int(resize_h)), (123.68, 116.78, 103.94), True, False)
    ratio_h = h / float(resize_h)
    ratio_w = w / float(resize_w)

    return blob, (ratio_h, ratio_w)


def mergeBoundingBox(cluster, boxes, indices, rW, rH):
    # bounding boxes points in counterclockwise order
    bounds = []
    for key in cluster:
        if key == -1:
            continue
        # unbound int :p
        min_x = min_y = 1000000
        max_x = max_y = -1000000
        box_num = ""
        for word in cluster[key]:
            box_num += f'{word} '
            vertices = cv.boxPoints(boxes[indices[word]])
            vertices = cv.boxPoints(boxes[indices[word]])
            for j in range(4):
                min_x = min(min_x, vertices[j][0])
                max_x = max(max_x, vertices[j][0])
                min_y = min(min_y, vertices[j][1])
                max_y = max(max_y, vertices[j][1])
        yeet = (
            min_x*rW, min_y*rH, max_x*rW, max_y*rH)
        (min_x, min_y, max_x, max_y) = [
            np.round(x).astype("int") for x in yeet]
        p1 = (min_x, min_y)
        p2 = (min_x, max_y)
        p3 = (max_x, max_y)
        p4 = (max_x, min_y)
        bounds.append((p1, p2, p3, p4))
    return bounds


def nothing(x):
    pass


def main():
    # Read and store arguments
    inPath = args.input
    outPath = args.output
    confThreshold = args.thr
    nmsThreshold = args.nms
    inpWidth = args.width
    inpHeight = args.height
    modelDetector = args.model    
    printReference = args.generate_reference
    scriptDir = os.path.dirname(__file__)
    outputList = [int(x) for x in args.list]
    if len(outputList) == 0 and not printReference:
        return



    detector = cv.dnn.readNet(scriptDir+ '/' + modelDetector)

    tickmeter = cv.TickMeter()

    frame = cv.imread(inPath)
    height_ = frame.shape[0]
    width_ = frame.shape[1]

    x_bias = 1
    y_bias = 10
    maxD = 0.2*width_
    print(maxD)
    alpha = 0.5

    inpHeight = height_ - (height_ % 32)
    inpWidth = width_ - (width_ % 32)

    rW = width_ / float(inpWidth)
    rH = height_ / float(inpHeight)
    blob, (rH, rW) = resize_image_blob(frame)

    tickmeter.start()
    outs = east(detector, blob, inpWidth, inpHeight)
    tickmeter.stop()
    # Get scores and geometry
    (scores, geometry) = outs
    (boxes, confidences) = decodeBoundingBoxes(
        scores, geometry, confThreshold)
    # Apply NMS
    indices = cv.dnn.NMSBoxesRotated(
        boxes, confidences, confThreshold, nmsThreshold)

    indices = indices.flatten()
    centers = []

    for i in indices:
        (center_x, center_y) = np.round(boxes[i][0]).astype("int")
        centers.append((center_x*x_bias, center_y*y_bias))
    print(centers)
    # Do camera stuff
    db = DBSCAN(eps=maxD, min_samples=1).fit(centers)

    cluster = {}

    for index, key in enumerate(db.labels_):
        if key in cluster:
            cluster[key].append(index)
        else:
            cluster[key] = [index]
    resetFrame = frame.copy()
    overlay = frame.copy()
    output = frame.copy()
    clusteredBound = frame.copy()

    boundingBoxes = mergeBoundingBox(cluster, boxes, indices, rW, rH)

    # Generate image for alpha
    for index, box in enumerate(boundingBoxes):
        if not printReference and index not in outputList:
            continue
        (p1, p2, p3, p4) = box
        cv.rectangle(overlay, p2, p4, (255, 0, 0), -1)
        cv.putText(overlay, str(index), p3,
                   cv.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), thickness=2)

    # Show Clustered bounding boxes
    # for index, box in enumerate(boundingBoxes):
    #     (p1, p2, p3, p4) = box
    #     cv.rectangle(clusteredBound, p2, p4, (255, 0, 0), 1)
    # Show word bounding boxes
    # for index, i in enumerate(indices):
    #     vertices = cv.boxPoints(boxes[i])
    #     # scale the bounding box coordinates based on the respective ratios
    #     for j in range(4):
    #         vertices[j][0] *= rW
    #         vertices[j][1] *= rH
    #     for j in range(4):
    #         p1 = (vertices[j][0], vertices[j][1])
    #         p2 = (vertices[(j + 1) % 4][0], vertices[(j + 1) % 4][1])
    #         cv.line(frame, p1, p2, (0, 255, 0), 1)

        # Put efficiency information
    label = 'Inference time: %.2f ms' % (tickmeter.getTimeMilli())
    cv.putText(frame, label, (0, 15),
               cv.FONT_HERSHEY_SIMPLEX, 0.5, (0, 255, 0))

    cv.addWeighted(overlay, alpha, output, 1 - alpha,
                   0, output)

    cv.imwrite(outPath, output)
    # cv.imwrite("test_box.jpg", clustered_bound)
    # cv.imwrite("test.jpg", frame)


if __name__ == "__main__":
    main()
