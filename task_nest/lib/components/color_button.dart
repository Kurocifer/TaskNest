import 'package:flutter/material.dart';
import 'package:task_nest/constants/constants.dart';

class ColorButton extends StatelessWidget {
  final void Function(int) changeColor;
  final ColorSelection colorSelected;

  const ColorButton({
    super.key,
    required this.changeColor,
    required this.colorSelected,
  });

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton(
      tooltip: 'change color',
      icon: Icon(
        Icons.opacity_outlined,
        color: Theme.of(context).colorScheme.onSurfaceVariant,
      ),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(10.0),
      ),
      itemBuilder: (context) {
        return List.generate(
          ColorSelection.values.length,
          (index) {
            final currentColor = ColorSelection.values[index];
            return PopupMenuItem(
              value: index,
              enabled: currentColor != colorSelected,
              child: Wrap(
                children: [
                  Padding(
                    padding: const EdgeInsets.only(left: 10),
                    child: Icon(
                      Icons.opacity_outlined,
                      color: currentColor.color,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.only(left: 20),
                    child: Text(currentColor.label),
                  ),
                ],
              ),
            );
          },
        );
      },
      onSelected: changeColor,
    );
  }
}
