import 'dart:io';

import 'package:flutter/material.dart';
import 'package:camera/camera.dart';
import 'package:image_picker/image_picker.dart';

class CameraGalleryScreen extends StatefulWidget {
  final List<CameraDescription> cameras;

  CameraGalleryScreen({required this.cameras});

  @override
  _CameraGalleryScreenState createState() => _CameraGalleryScreenState();
}

class _CameraGalleryScreenState extends State<CameraGalleryScreen> {
  CameraController? controller;
  String? imagePath;
  final ImagePicker _picker = ImagePicker();

  @override
  void initState() {
    super.initState();
    initializeCamera();
  }

  Future<void> initializeCamera() async {
    if (widget.cameras.isNotEmpty) {
      controller = CameraController(widget.cameras[0], ResolutionPreset.high);
      await controller?.initialize();
      setState(() {});
    }
  }

  @override
  void dispose() {
    controller?.dispose();
    super.dispose();
  }

  Future<void> takePicture() async {
    if (!controller!.value.isInitialized) {
      return;
    }
    if (controller!.value.isTakingPicture) {
      return;
    }

    try {
      XFile picture = await controller!.takePicture();
      setState(() {
        imagePath = picture.path;
      });
    } catch (e) {
      print(e);
    }
  }

  Future<void> pickImageFromGallery() async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      setState(() {
        imagePath = image.path;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (controller == null || !controller!.value.isInitialized) {
      return Container();
    }
    return Scaffold(
      appBar: AppBar(
        title: Text('Camera and Gallery'),
      ),
      body: Column(
        children: [
          if (imagePath == null)
            Expanded(
              child: CameraPreview(controller!),
            )
          else
            Expanded(
              child: Image.file(File(imagePath!)),
            ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                ElevatedButton(
                  onPressed: takePicture,
                  child: Text('Take Picture'),
                ),
                ElevatedButton(
                  onPressed: pickImageFromGallery,
                  child: Text('Load from Gallery'),
                ),
                if (imagePath != null)
                  ElevatedButton(
                    onPressed: () {
                      setState(() {
                        imagePath = null;
                      });
                    },
                    child: Text('Retake/Reload'),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
