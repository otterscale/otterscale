import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter/semantics.dart';
import 'package:ui/home/home.dart';
import 'package:ui/notifier/notifier.dart';
import 'package:watch_it/watch_it.dart';
import 'package:yaru/yaru.dart';

Future<void> main() async {
  await YaruWindowTitleBar.ensureInitialized();
  WidgetsFlutterBinding.ensureInitialized();
  SemanticsBinding.instance.ensureSemantics();

  registerManager();
  runApp(const MyApp());
}

void registerManager() {
  di.registerSingleton<ThemeModeManager>(ThemeModeManager());
}

class MyApp extends StatelessWidget with WatchItMixin {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    final themeMode = watchPropertyValue((ThemeModeManager m) => m.themeMode);

    return YaruTheme(
      builder: (context, value, child) => MaterialApp(
        home: const MyHomePage(title: 'Flutter Demo Home Page'),
        title: 'OpenHDC',
        theme: ThemeData.light(useMaterial3: true),
        darkTheme: ThemeData.dark(useMaterial3: true),
        themeMode: themeMode,
        debugShowCheckedModeBanner: false,
        scrollBehavior: const MaterialScrollBehavior().copyWith(
          dragDevices: {
            PointerDeviceKind.mouse,
            PointerDeviceKind.touch,
            PointerDeviceKind.stylus,
            PointerDeviceKind.unknown,
            PointerDeviceKind.trackpad,
          },
        ),
      ),
    );
  }
}
