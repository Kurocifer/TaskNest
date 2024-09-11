import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import './screens/home.dart';

void main() {
  runApp(const TaskNest());
}

class TaskNest extends StatefulWidget {
  const TaskNest({super.key});

  @override
  State<TaskNest> createState() => _TaskNestState();
}

class _TaskNestState extends State<TaskNest> {
  // This widget is the root of your application.
  ThemeMode themeMode = ThemeMode.light; // default theme
  Color colorSelected = Colors.white; // default app color

  void changeThemeMode(bool useLightMode) {
    setState(() {
      themeMode = useLightMode ? ThemeMode.light : ThemeMode.dark;
    });
  }

  @override
  Widget build(BuildContext context) {
    SystemChrome.setSystemUIOverlayStyle(
        const SystemUiOverlayStyle(statusBarColor: Colors.transparent));

    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'ToDo App',
      themeMode: themeMode,
      theme: ThemeData(
        colorSchemeSeed: colorSelected,
        useMaterial3: true,
        brightness: Brightness.light,
      ),
      darkTheme: ThemeData(
        colorSchemeSeed: colorSelected,
        useMaterial3: true,
        brightness: Brightness.dark,
      ),
      home: Home(
        changeTheme: changeThemeMode,
        colorSelected: colorSelected,
        isLightMode: themeMode == ThemeMode.light,
      ),
    );
  }
}
