# Instance Scheduler

A service that is responsible for powering on and off compute instances based on a schedule. This
schedule is defined as a tag on the cloud instance resource.

## Tags

- **InstanceSchedulerEnabled** (`bool`)
- **InstanceSchedulerSchedule** (`json`)
  ```json
  {
    "default": "09:00-17:00",
    "overrides": {
      "monday": ["09:00-12:00", "17:00-21:00"]
      "saturday": ["-"],
      "sunday": ["-"]
    }
  }
  ```
The schedule is comprised of the `default` and `overrides` keys. The `default` key is a string that
has a time range that the virtual machine will be powered on for, all other time it will be powered
off. The `overrides` key is a map where the key states which day will be overridden and the value is
a slice of either time ranges or a `-` which signifies that the machine will be off entirely for that
day.
