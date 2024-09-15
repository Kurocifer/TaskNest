import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

/// Credential Class
class Credentials {
  Credentials(this.username, this.password);
  final String username;
  final String password;
}

class LoginPage extends StatelessWidget {
  /// Called when users sign in with [Credentials].
  final ValueChanged<Credentials> onLogIn;

  const LoginPage({required this.onLogIn, super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: LayoutBuilder(
        builder: (BuildContext context, BoxConstraints constraints) {
          if (constraints.maxWidth > 700) {
            // Display Desktop View
            return Row(
              children: [
                Expanded(
                  child: SizedBox(
                    width: double.infinity,
                    height: double.infinity,
                    child: Image.asset(
                      'assets/images/login_background.jpg',
                      fit: BoxFit.cover,
                    ),
                  ),
                ),
                Expanded(
                  child: FractionallySizedBox(
                      widthFactor: 0.70,
                      child: LoginForm(
                        onLogIn: onLogIn,
                      )),
                ),
              ],
            );
          } else {
            // Display Mobile View
            return Column(
              children: [
                Expanded(
                  child: LoginForm(onLogIn: onLogIn),
                ),
              ],
            );
          }
        },
      ),
    );
  }
}

class LoginForm extends StatelessWidget {
  LoginForm({required this.onLogIn, super.key});

  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final ValueChanged<Credentials> onLogIn;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          SizedBox(
            height: 40.0,
            width: 40.0,
            child: ClipRRect(
                borderRadius: BorderRadius.circular(20.0),
                child: Image.asset('assets/images/cifer.png')),
          ),
          const SizedBox(height: 20),
          TextField(
            controller: _usernameController,
            decoration: InputDecoration(
              // filled: true,
              hintText: 'Username',
              border: OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(8.0), // Rounded corner border
              ),
            ),
          ),
          const SizedBox(height: 12),
          TextField(
            controller: _passwordController,
            obscureText: true,
            decoration: InputDecoration(
              // filled: true,
              hintText: 'Password',
              border: OutlineInputBorder(
                borderRadius:
                    BorderRadius.circular(8.0), // Rounded corner border
              ),
            ),
          ),
          const SizedBox(height: 24),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              ElevatedButton(
                child: const Text('Login'),
                onPressed: () async {
                  onLogIn(Credentials(_usernameController.value.text,
                      _passwordController.value.text));
                },
              ),
              ElevatedButton(
                child: const Text("Don't Yet have an Account"),
                onPressed: () {
                  context.go('/signup');
                },
              ),
            ],
          )
        ],
      ),
    );
  }
}
