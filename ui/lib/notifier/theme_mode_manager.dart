import 'package:flutter/material.dart';
import 'package:safe_change_notifier/safe_change_notifier.dart';

class ThemeModeManager extends SafeChangeNotifier {
  ThemeModeManager();

  ThemeMode _themeMode = ThemeMode.system;
  ThemeMode get themeMode => _themeMode;

  void setThemeMode(ThemeMode value) {
    if (value == _themeMode) return;
    _themeMode = value;
    notifyListeners();
  }
}
