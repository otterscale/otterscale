import 'package:flutter/material.dart';
import 'package:ui/notifier/notifier.dart';
import 'package:watch_it/watch_it.dart';

class ThemeModeButton extends StatelessWidget {
  const ThemeModeButton({super.key});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return IconButton(
      onPressed: () => di<ThemeModeManager>()
          .setThemeMode(isDark ? ThemeMode.light : ThemeMode.dark),
      icon:
          Icon(isDark ? Icons.wb_sunny_outlined : Icons.brightness_2_outlined),
    );
  }
}
