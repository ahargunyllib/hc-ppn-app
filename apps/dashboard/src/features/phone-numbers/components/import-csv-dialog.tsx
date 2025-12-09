import { Button } from "@/shared/components/ui/button";
import {
  Dialog,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogPanel,
  DialogPopup,
  DialogTitle,
  DialogTrigger,
} from "@/shared/components/ui/dialog";
import { Input } from "@/shared/components/ui/input";
import { toastManager } from "@/shared/components/ui/toast";
import { parseAPIError } from "@/shared/lib/api-client";
import { useImportUsersFromCSV } from "@/shared/repositories/user/query";
import { Upload } from "lucide-react";
import { useRef, useState } from "react";

export default function ImportCSVDialog() {
  const [file, setFile] = useState<File | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const { mutate, isPending } = useImportUsersFromCSV();
  const inputRef = useRef<HTMLInputElement | null>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const f = event.target.files?.[0];
    if (!f) {
      return;
    }

    setFile(f);
  };

  const onSubmitHandler = () => {
    if (!file) {
      toastManager.add({
        type: "error",
        title: "No file selected",
        description: "Please select a CSV file to import.",
      });
      return;
    }

    mutate(file, {
      onSuccess: () => {
        toastManager.add({
          type: "success",
          title: "CSV imported successfully",
        });
      },
      onError: (error) => {
        toastManager.add({
          type: "error",
          title: "Gagal mengimpor CSV",
          description: parseAPIError(error),
        });
      },
    });
  };

  return (
    <Dialog onOpenChange={setIsOpen} open={isOpen}>
      <DialogTrigger render={<Button size="sm" variant="outline" />}>
        <Upload />
        Import CSV
      </DialogTrigger>
      <DialogPopup>
        <DialogHeader>
          <DialogTitle>Import Users dari CSV</DialogTitle>
          <DialogDescription>
            Upload file CSV untuk menambahkan multiple users sekaligus.
          </DialogDescription>
        </DialogHeader>
        <DialogPanel>
          <div className="flex flex-col gap-4">
            <div className="space-y-2 rounded-md border border-border bg-muted p-4">
              <h3 className="font-semibold text-sm">Format CSV:</h3>
              <div className="space-y-2 text-sm">
                <p>File CSV harus memiliki kolom berikut (dengan header):</p>
                <ul className="ml-4 list-disc space-y-1">
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      phoneNumber
                    </code>{" "}
                    (wajib) - Format E.164, contoh: +6281234567890
                  </li>
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      name
                    </code>{" "}
                    (wajib) - Nama lengkap
                  </li>
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      jobTitle
                    </code>{" "}
                    (opsional) - Jabatan
                  </li>
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      gender
                    </code>{" "}
                    (opsional) - male atau female
                  </li>
                  <li>
                    <code className="rounded bg-background px-1 py-0.5">
                      dateOfBirth
                    </code>{" "}
                    (opsional) - Format: YYYY-MM-DD
                  </li>
                </ul>
                <div className="mt-3">
                  <p className="font-semibold">Contoh:</p>
                  <pre className="mt-1 overflow-x-auto rounded bg-background p-2 font-mono text-xs">
                    {`phoneNumber,name,jobTitle,gender,dateOfBirth
+6281234567890,John Doe,Software Engineer,male,1990-01-15
+6281234567891,Jane Smith,Product Manager,female,1992-05-20`}
                  </pre>
                </div>
              </div>
            </div>

            <div>
              <Input
                accept=".csv"
                disabled={isPending}
                onChange={handleFileChange}
                ref={inputRef}
                type="file"
              />
            </div>
          </div>
        </DialogPanel>
        <DialogFooter>
          <Button
            disabled={isPending}
            onClick={() => setIsOpen(false)}
            variant="outline"
          >
            {isPending ? "Mengimpor..." : "Tutup"}
          </Button>
          <Button disabled={isPending} onClick={onSubmitHandler}>
            {isPending ? "Mengimpor..." : "Import"}
          </Button>
        </DialogFooter>
      </DialogPopup>
    </Dialog>
  );
}
