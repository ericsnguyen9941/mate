import AppKit
import UserNotifications

struct Message: Decodable {
    let id: String
    let title: String
    let body: String
    let url: String?
}

final class Delegate: NSObject, NSApplicationDelegate, UNUserNotificationCenterDelegate {
    var resultPath: String?
    var finished = false

    func finish(_ error: String = "") {
        DispatchQueue.main.async {
            guard !self.finished else { return }
            self.finished = true
            if let path = self.resultPath,
               let data = try? JSONSerialization.data(withJSONObject: ["error": error]) {
                try? data.write(to: URL(fileURLWithPath: path), options: .atomic)
            }
            NSApplication.shared.terminate(nil)
        }
    }

    func applicationDidFinishLaunching(_ notification: Notification) {
        let center = UNUserNotificationCenter.current()
        center.delegate = self
        let args = CommandLine.arguments
        guard args.count == 4, args[1] == "--send" else {
            // Notification Center may launch the app again when a delivered
            // notification is clicked, after its sending process has exited.
            DispatchQueue.main.asyncAfter(deadline: .now() + 15) { self.finish() }
            return
        }
        resultPath = args[3]
        DispatchQueue.main.asyncAfter(deadline: .now() + 90) {
            self.finish("Timed out waiting for macOS notification permission; enable Mate in System Settings > Notifications.")
        }
        do {
            let message = try JSONDecoder().decode(Message.self, from: Data(contentsOf: URL(fileURLWithPath: args[2])))
            center.requestAuthorization(options: [.alert, .sound]) { allowed, error in
                guard allowed, error == nil else {
                    self.finish(error?.localizedDescription ?? "Mate notifications are disabled. Enable Mate in System Settings > Notifications.")
                    return
                }
                let content = UNMutableNotificationContent()
                content.title = message.title
                content.body = message.body
                content.sound = .default
                if let url = message.url { content.userInfo = ["url": url] }
                let request = UNNotificationRequest(identifier: message.id, content: content, trigger: nil)
                center.add(request) { error in
                    // The receipt means accepted by Notification Center, not
                    // proof a banner was displayed (Focus/settings control that).
                    self.finish(error?.localizedDescription ?? "")
                }
            }
        } catch { finish(error.localizedDescription) }
    }

    func userNotificationCenter(_ center: UNUserNotificationCenter,
                                willPresent notification: UNNotification,
                                withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions) -> Void) {
        completionHandler([.banner, .list, .sound])
    }

    func userNotificationCenter(_ center: UNUserNotificationCenter,
                                didReceive response: UNNotificationResponse,
                                withCompletionHandler completionHandler: @escaping () -> Void) {
        if let value = response.notification.request.content.userInfo["url"] as? String,
           let url = URL(string: value), url.scheme == "http",
           ["127.0.0.1", "localhost", "[::1]"].contains(url.host ?? "") {
            NSWorkspace.shared.open(url)
        }
        completionHandler()
        finish()
    }
}

let app = NSApplication.shared
let delegate = Delegate()
app.delegate = delegate
app.setActivationPolicy(.accessory)
app.run()
