import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'package:task_nest/constants/constants.dart';
import 'package:task_nest/model/auth.dart';
import 'package:task_nest/screens/login_page.dart';
import 'package:task_nest/screens/signup_page.dart';
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
  ColorSelection colorSelected = ColorSelection.white; // default app color
  final TaskNestAuth _auth = TaskNestAuth();

  late final _router = GoRouter(
      initialLocation: '/login',
      redirect: _appRedirect,
      routes: [
        GoRoute(
          path: '/signup',
          builder: (context, state) => SignUpPage(
            onSignUP: (SignUpCredentials credentials) async {
              _auth
                  .signUp(credentials.username, credentials.password)
                  .then((_) => context.go('/${TaskNestTab.home.value}'));
            },
          ),
        ),
        GoRoute(
          path: '/login',
          builder: (context, state) => LoginPage(
            onLogIn: (Credentials credentials) async {
              _auth
                  .signIn(credentials.username, credentials.password)
                  .then((_) => context.go('/${TaskNestTab.home.value}'));
            },
          ),
        ),
        GoRoute(
          path: '/:tab',
          builder: (context, state) {
            return Home(
              auth: _auth,
              changeTheme: changeThemeMode,
              changeColor: changeColor,
              colorSelected: colorSelected,
              tab: 0,
            );
          },
        ),
      ],
      errorPageBuilder: (context, state) {
        return MaterialPage(
          key: state.pageKey,
          child: Scaffold(
            appBar: AppBar(
              leading: IconButton(
                icon: const Icon(Icons.arrow_back),
                onPressed: () => Navigator.of(context).pop(),
              ),
            ),
            body: Center(
              child: Text(
                state.error.toString(),
              ),
            ),
          ),
        );
      });

  Future<String?> _appRedirect(
      BuildContext context, GoRouterState state) async {
    final loggedIn = await _auth.loggedIn;
    final isOnLoginPage = state.matchedLocation == '/login';

    // Go to /login if the user is not signed in
    //if (!loggedIn) {
    //  return '/login';
    //}

    // Go to root if the user is already signed in
    if (loggedIn && isOnLoginPage) {
      return '/${TaskNestTab.home.value}';
    }

    
    return null;
  }

  void changeColor(int value) {
    setState(() {
      colorSelected = ColorSelection.values[value];
    });
  }

  void changeThemeMode(bool useLightMode) {
    setState(() {
      themeMode = useLightMode ? ThemeMode.light : ThemeMode.dark;
    });
  }

  @override
  Widget build(BuildContext context) {
    SystemChrome.setSystemUIOverlayStyle(
        const SystemUiOverlayStyle(statusBarColor: Colors.transparent));

    return MaterialApp.router(
      debugShowCheckedModeBanner: false,
      title: 'ToDo App',
      themeMode: themeMode,
      routerConfig: _router,
      theme: ThemeData(
        colorSchemeSeed: colorSelected.color,
        useMaterial3: true,
        brightness: Brightness.light,
      ),
      darkTheme: ThemeData(
        colorSchemeSeed: colorSelected.color,
        useMaterial3: true,
        brightness: Brightness.dark,
      ),
      /*
      home: Home(
        changeTheme: changeThemeMode,
        changeColor: changeColor,
        colorSelected: colorSelected,
      ),*/
    );
  }
}
