// Go-Reloaded GitHub Pages Interactive Features

// FSM State Animation
function animateFSM() {
    const textState = document.getElementById('state-text');
    const commandState = document.getElementById('state-command');
    
    setInterval(() => {
        if (textState.classList.contains('active')) {
            textState.classList.remove('active');
            commandState.classList.add('active');
        } else {
            commandState.classList.remove('active');
            textState.classList.add('active');
        }
    }, 2000);
}

// Text Transformation Demo
function initDemo() {
    const transformBtn = document.getElementById('transform-btn');
    const inputText = document.getElementById('demo-input-text');
    const outputText = document.getElementById('demo-output-text');
    
    if (!transformBtn || !inputText || !outputText) return;
    
    transformBtn.addEventListener('click', () => {
        const input = inputText.value;
        const output = simulateTransformation(input);
        outputText.textContent = output;
        
        // Add animation
        outputText.style.opacity = '0';
        setTimeout(() => {
            outputText.style.opacity = '1';
        }, 100);
    });
}

// Simulate Go-Reloaded transformations (client-side approximation)
function simulateTransformation(text) {
    let result = text;
    
    // Hex conversions
    result = result.replace(/([A-Fa-f0-9]+)\s*\(hex\)/g, (match, hex) => {
        try {
            return parseInt(hex, 16).toString();
        } catch {
            return match;
        }
    });
    
    // Binary conversions
    result = result.replace(/([01]+)\s*\(bin\)/g, (match, bin) => {
        try {
            return parseInt(bin, 2).toString();
        } catch {
            return match;
        }
    });
    
    // Case transformations - single word
    result = result.replace(/(\w+)\s*\(up\)/g, (match, word) => word.toUpperCase());
    result = result.replace(/(\w+)\s*\(low\)/g, (match, word) => word.toLowerCase());
    result = result.replace(/(\w+)\s*\(cap\)/g, (match, word) => 
        word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()
    );
    
    // Multi-word transformations (simplified)
    result = result.replace(/(\w+(?:\s+\w+)*)\s*\(up,\s*(\d+)\)/g, (match, words, count) => {
        const wordArray = words.split(/\s+/);
        const numWords = Math.min(parseInt(count), wordArray.length);
        for (let i = wordArray.length - numWords; i < wordArray.length; i++) {
            wordArray[i] = wordArray[i].toUpperCase();
        }
        return wordArray.join(' ');
    });
    
    result = result.replace(/(\w+(?:\s+\w+)*)\s*\(cap,\s*(\d+)\)/g, (match, words, count) => {
        const wordArray = words.split(/\s+/);
        const numWords = Math.min(parseInt(count), wordArray.length);
        for (let i = wordArray.length - numWords; i < wordArray.length; i++) {
            wordArray[i] = wordArray[i].charAt(0).toUpperCase() + wordArray[i].slice(1).toLowerCase();
        }
        return wordArray.join(' ');
    });
    
    result = result.replace(/(\w+(?:\s+\w+)*)\s*\(low,\s*(\d+)\)/g, (match, words, count) => {
        const wordArray = words.split(/\s+/);
        const numWords = Math.min(parseInt(count), wordArray.length);
        for (let i = wordArray.length - numWords; i < wordArray.length; i++) {
            wordArray[i] = wordArray[i].toLowerCase();
        }
        return wordArray.join(' ');
    });
    
    // Article corrections
    result = result.replace(/\ba\s+([aeiouAEIOU])/g, 'an $1');
    result = result.replace(/\ban\s+([bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ])/g, 'a $1');
    result = result.replace(/\ba\s+h/g, 'an h'); // Special case for 'h'
    
    // Punctuation spacing
    result = result.replace(/\s+([,.!?;:])/g, '$1');
    result = result.replace(/([,.!?;:])(?!\s)/g, '$1 ');
    
    // Clean up extra spaces
    result = result.replace(/\s+/g, ' ').trim();
    
    return result;
}

// Smooth scrolling for navigation links
function initSmoothScrolling() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });
}

// Navbar scroll effect
function initNavbarScroll() {
    const navbar = document.querySelector('.navbar');
    if (!navbar) return;
    
    window.addEventListener('scroll', () => {
        if (window.scrollY > 100) {
            navbar.style.background = 'rgba(255, 255, 255, 0.98)';
            navbar.style.boxShadow = '0 4px 6px -1px rgba(0, 0, 0, 0.1)';
        } else {
            navbar.style.background = 'rgba(255, 255, 255, 0.95)';
            navbar.style.boxShadow = 'none';
        }
    });
}

// Intersection Observer for animations
function initScrollAnimations() {
    const observerOptions = {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
    };
    
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.style.opacity = '1';
                entry.target.style.transform = 'translateY(0)';
            }
        });
    }, observerOptions);
    
    // Observe feature cards
    document.querySelectorAll('.feature-card, .example-card, .fsm-detail').forEach(card => {
        card.style.opacity = '0';
        card.style.transform = 'translateY(30px)';
        card.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
        observer.observe(card);
    });
}

// Copy code examples
function initCodeCopy() {
    document.querySelectorAll('.feature-example code, .example-demo code').forEach(codeBlock => {
        codeBlock.style.cursor = 'pointer';
        codeBlock.title = 'Click to copy';
        
        codeBlock.addEventListener('click', () => {
            navigator.clipboard.writeText(codeBlock.textContent).then(() => {
                const originalText = codeBlock.textContent;
                codeBlock.textContent = 'Copied!';
                codeBlock.style.background = '#10b981';
                codeBlock.style.color = 'white';
                
                setTimeout(() => {
                    codeBlock.textContent = originalText;
                    codeBlock.style.background = '#f3f4f6';
                    codeBlock.style.color = 'inherit';
                }, 1000);
            });
        });
    });
}

// Performance chart animation
function animatePerformanceChart() {
    const chartPoints = document.querySelectorAll('.chart-point');
    chartPoints.forEach((point, index) => {
        setTimeout(() => {
            point.style.opacity = '1';
            point.style.transform = 'translateX(-50%) translateY(0)';
        }, index * 200);
    });
}

// Initialize all features when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    animateFSM();
    initDemo();
    initSmoothScrolling();
    initNavbarScroll();
    initScrollAnimations();
    initCodeCopy();
    
    // Animate performance chart when it comes into view
    const performanceSection = document.getElementById('performance');
    if (performanceSection) {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    animatePerformanceChart();
                    observer.unobserve(entry.target);
                }
            });
        });
        observer.observe(performanceSection);
    }
    
    // Add loading animation
    document.body.style.opacity = '0';
    setTimeout(() => {
        document.body.style.transition = 'opacity 0.5s ease';
        document.body.style.opacity = '1';
    }, 100);
});

// Add some interactive examples
const interactiveExamples = [
    {
        input: "The value is FF (hex)",
        output: "The value is 255"
    },
    {
        input: "Binary 1010 (bin) equals decimal",
        output: "Binary 10 equals decimal"
    },
    {
        input: "these three words (up, 3)",
        output: "THESE THREE WORDS"
    },
    {
        input: "I need a apple and an car",
        output: "I need an apple and a car"
    },
    {
        input: "Hello , world ! How are you ?",
        output: "Hello, world! How are you?"
    }
];

// Cycle through examples in demo
function cycleExamples() {
    const inputTextarea = document.getElementById('demo-input-text');
    if (!inputTextarea) return;
    
    let currentExample = 0;
    
    setInterval(() => {
        if (document.activeElement !== inputTextarea) {
            inputTextarea.value = interactiveExamples[currentExample].input;
            currentExample = (currentExample + 1) % interactiveExamples.length;
        }
    }, 5000);
}

// Start example cycling after page load
setTimeout(cycleExamples, 3000);