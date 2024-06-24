import 'dart:io';

import 'package:http/http.dart' as http;
import 'package:path_provider/path_provider.dart';
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
  final ImagePicker _picker = ImagePicker();

  File? _image;
  File? _pImage;

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
        _image = File(picture.path);
      });
    } catch (e) {
      print(e);
    }
  }

  Future<void> pickImageFromGallery() async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      setState(() {
        _image = File(image.path);
      });
    }
  }

  Future<void> _processImage() async {
    if (_image == null) return;

    final request = http.MultipartRequest(
      'POST',
      Uri.parse('http://10.0.0.97:8080/process_image'),
    );
    request.files.add(await http.MultipartFile.fromPath('image', _image!.path));

    final response = await request.send();

    if (response.statusCode == 200) {
      final Directory tempDir = await getTemporaryDirectory();
      final String uniqueFileName =
          'processed_image_${DateTime.now().millisecondsSinceEpoch}.png';
      final File file = File('${tempDir.path}/$uniqueFileName');
      await response.stream.pipe(file.openWrite());
      setState(() {
        _pImage = file;
      });
    } else {
      print('Failed to process image');
    }
  }

  Widget _mainDisplay() {
    if (_pImage != null) {
      return Expanded(child: Image.file(_pImage!));
    }

    if (_image == null) {
      return Expanded(child: CameraPreview(controller!));
    }

    return Expanded(child: Image.file(_image!));
  }

  @override
  Widget build(BuildContext context) {
    if (controller == null || !controller!.value.isInitialized) {
      return Container();
    }
    return Scaffold(
      appBar: AppBar(
        title: const Text('Camera and Gallery'),
      ),
      body: Column(
        children: [
          _mainDisplay(),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                if (_image == null)
                  ElevatedButton(
                    onPressed: takePicture,
                    child: const Text('Take Picture'),
                  ),
                if (_image == null)
                  ElevatedButton(
                    onPressed: pickImageFromGallery,
                    child: const Text('Load from Gallery'),
                  ),
                if (_image != null)
                  ElevatedButton(
                    onPressed: () {
                      setState(() {
                        _image = null;
                        _pImage = null;
                      });
                    },
                    child: const Text('Retake/Reload'),
                  ),
                if (_image != null)
                  ElevatedButton(
                    onPressed: _processImage,
                    child: const Text('Process'),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
