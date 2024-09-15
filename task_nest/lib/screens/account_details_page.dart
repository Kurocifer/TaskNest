import 'package:flutter/material.dart';

class AccountDetailsPage extends StatefulWidget {
  final Function onLogOut;

  const AccountDetailsPage({
    super.key,
    required this.onLogOut,
  });

  @override
  State<AccountDetailsPage> createState() => _AccountDetailsPageState();
}

class _AccountDetailsPageState extends State<AccountDetailsPage> {
  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context)
        .textTheme
        .apply(displayColor: Theme.of(context).colorScheme.onSurface);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              'Your Details',
              style: textTheme.headlineSmall,
            ),
            Center(
              child: ListTile(
                title: const Text(
                  'Log out',
                  style: TextStyle(fontSize: 30),
                ),
                onTap: () {
                  widget.onLogOut(true);
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
