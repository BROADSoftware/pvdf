
# TODO

- If the pvscanner daemonset does not run, whatever the reason is, the usage values are wrong.
  A solution could be to check pvscanner good health in pvdf. (Ensure they are in a running state)

- Also, test and issue a warning on each shot when vgsd is unreachable

- A hidden cleanup command to clean all annotation

- The way Udpate usage information for Node 's1':size.topolvm.pvscanner.pvdf.broadsoftware.com/nvme:107369988096 is displayed is wrong. Issue this message in adjustAnnotation()