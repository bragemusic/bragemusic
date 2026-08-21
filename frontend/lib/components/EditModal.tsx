import { Button, Form, useOverlayState } from "@heroui/react";
import { Modal } from "@heroui/react";
import { Edit, LucideIcon, Save } from "lucide-react";

type State = ReturnType<typeof useOverlayState>;

export default function EditModal({
  state,
  children,
  header,
  onSubmit,
  icon: Icon,
  canSubmit = true,
}: {
  state: State;
  children: React.ReactNode;
  header: string;
  onSubmit: (data: FormData) => void;
  icon?: LucideIcon;
  canSubmit?: boolean;
}) {
  return (
    <Modal>
      <Modal.Backdrop
        isOpen={state.isOpen}
        variant="opaque"
        onOpenChange={state.setOpen}
      >
        <Modal.Container size="lg">
          <Modal.Dialog className="">
            <Modal.CloseTrigger />
            <Modal.Header>
              <Modal.Icon className="bg-accent-soft text-accent-soft-foreground">
                {Icon ? <Icon /> : <Edit />}
              </Modal.Icon>
              <Modal.Heading className="font-semibold">{header}</Modal.Heading>
            </Modal.Header>
            <Modal.Body className="px-1 pt-4">
              <Form
                className="flex flex-col gap-4"
                id="modalform"
                onSubmit={(e) => {
                  e.preventDefault();
                  onSubmit(new FormData(e.currentTarget));
                }}
              >
                {children}
              </Form>
            </Modal.Body>
            <Modal.Footer>
              <Button slot="close" variant="secondary">
                Cancel
              </Button>
              {canSubmit &&
                <Button
                  form="modalform"
                  slot="close"
                  type="submit"
                  variant="primary"
                >
                  <Save />
                  Save
                </Button>
              }
            </Modal.Footer>
          </Modal.Dialog>
        </Modal.Container>
      </Modal.Backdrop>
    </Modal>
  );
}
