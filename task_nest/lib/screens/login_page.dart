import 'package:flutter/material.dart';

/// Credential Class
class Credentials {
  Credentials(this.username, this.password);
  final String username;
  final String password;
}

class LoginPage extends StatelessWidget {
  const LoginPage({required this.onLogIn, super.key});

  /// Called when users sign in with [Credentials].
  final ValueChanged<Credentials> onLogIn;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: LayoutBuilder(
        builder: (BuildContext context, BoxConstraints constraints) {
            return Column(
              children: [
                Expanded(
                  child: LoginForm(onLogIn: onLogIn),
                ),
              ],
            );
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
          ElevatedButton(
            child: const Text('Login'),
            onPressed: () async {
              onLogIn(Credentials(_usernameController.value.text,
                  _passwordController.value.text));
            },
          ),
        ],
      ),
    );
  }
}
