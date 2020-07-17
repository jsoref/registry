import 'package:flutter/material.dart';

import 'application.dart';
import 'theme.dart';
import 'routes.dart';
import 'signin.dart';

void main() async {
  await attemptToSignIn();
  runApp(Application());
}

class Application extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: applicationName,
      theme: appTheme(),
      initialRoute: "/",
      onGenerateRoute: generateRoute,
    );
  }
}
