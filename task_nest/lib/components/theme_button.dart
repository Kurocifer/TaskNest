import 'package:flutter/material.dart';

class ThemeButton extends StatelessWidget {
  final Function changeThemeMode;

  const ThemeButton({
    super.key,
    required this.changeThemeMode,
  });

  @override
  Widget build(BuildContext context) {
    final isBriget = Theme.of(context).brightness == Brightness.light;

    return IconButton(
      tooltip: isBriget ? 'use dark theme' : 'use ligt theme',
      icon: isBriget
          ? const Icon(Icons.dark_mode_outlined)
          : const Icon(Icons.light_mode_outlined),
      onPressed: () => changeThemeMode(!isBriget), 
    );
  }
}
