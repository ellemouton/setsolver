import 'package:flutter/material.dart';
import 'home_page.dart';

void main() {
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
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: const HomePage(title: 'SET Solver'),
    );
  }
}
