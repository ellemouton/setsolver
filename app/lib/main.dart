import 'package:flutter/material.dart';
import 'home_screen.dart';
import 'package:camera/camera.dart';

List<CameraDescription> cameras = [];

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  try {
    cameras = await availableCameras();
  } catch (e) {
    print('Error: $e');
  }

  runApp(const SetSolver());
}

// SetSolver is the main app.
class SetSolver extends StatelessWidget {
  // SetSolver constructor.
  const SetSolver({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'SET Solver',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
            seedColor: const Color.fromARGB(255, 100, 33, 217)),
        useMaterial3: true,
      ),
      home: const HomePage(title: 'SET Solver'),
    );
  }
}
